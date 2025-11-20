package tests

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"testing"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
	"github.com/google/uuid"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/cassandra"
	"github.com/testcontainers/testcontainers-go/modules/kafka"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/twmb/franz-go/pkg/kgo"
)

var (
	cassandraDB  CassandraTestDB
	kafkaCluster KafkaTestCluster
)

type CassandraTestDB struct {
	cluster   *gocql.ClusterConfig
	container *cassandra.CassandraContainer
}

type KafkaTestCluster struct {
	container *kafka.KafkaContainer
	brokers   []string
}

func (k KafkaTestCluster) NewTopicClient(t *testing.T) (*kgo.Client, string) {
	t.Helper()

	if k.container == nil {
		t.Fatal("kafka cluster not initialized")
	}

	ctx := t.Context()
	topic := "inventory_" + strings.ReplaceAll(uuid.NewString(), "-", "")
	if err := k.createTopic(ctx, topic); err != nil {
		t.Fatalf("failed to create topic %q: %v", topic, err)
	}

	client := must(kgo.NewClient(
		kgo.SeedBrokers(k.brokers...),
		kgo.ConsumerGroup("inventory-service"),
		kgo.ConsumeTopics(topic),
	))

	return client, topic
}

func (k KafkaTestCluster) createTopic(ctx context.Context, topic string) error {
	_, _, err := k.container.Exec(ctx, []string{
		"/usr/bin/kafka-topics",
		"--create",
		"--topic", topic,
		"--bootstrap-server", "localhost:9092",
		"--partitions", "1",
		"--replication-factor", "1",
	})
	return err
}

func (c CassandraTestDB) Session(t *testing.T) (*gocql.Session, error) {
	ctx := t.Context()

	// connect without keyspace
	root := *c.cluster
	root.Keyspace = ""
	rootSession, err := root.CreateSession()
	if err != nil {
		return nil, err
	}
	defer rootSession.Close()

	// create unique keyspace
	ks := "test_" + strings.ReplaceAll(uuid.NewString(), "-", "")
	if err := rootSession.Query(fmt.Sprintf(`
		CREATE KEYSPACE %s
		WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1}`, ks)).
		ExecContext(ctx); err != nil {
		return nil, err
	}

	// connect bound to that keyspace
	cluster := *c.cluster
	cluster.Keyspace = ks
	session, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}

	if err := schemaMigration(session); err != nil {
		return nil, err
	}

	return session, nil
}

func schemaMigration(session *gocql.Session) error {
	files, err := filepath.Glob("../../migrations/*cql")
	if err != nil {
		return err
	}

	slices.Sort(files)

	for _, f := range files {
		content, err := os.ReadFile(f)
		if err != nil {
			return err
		}
		for stmt := range strings.SplitSeq(string(content), ";") {
			if len(stmt) < 2 {
				continue
			}
			if err := session.Query(stmt).Exec(); err != nil {
				return err
			}

		}
	}

	return nil
}

func TestMain(m *testing.M) {
	var wg sync.WaitGroup
	ctx := context.Background()

	wg.Go(func() {
		waitLog := wait.ForLog("Starting listening for CQL clients")
		cassandraContainer := must(
			cassandra.Run(ctx, "cassandra:4.1.3", testcontainers.WithWaitStrategy(waitLog)),
		)
		host := must(cassandraContainer.Host(ctx))
		port := must(cassandraContainer.MappedPort(ctx, "9042/tcp"))
		cassandraDB = CassandraTestDB{
			cluster:   gocql.NewCluster(fmt.Sprintf("%s:%s", host, port.Port())),
			container: cassandraContainer,
		}
	})

	var kafkaContainer *kafka.KafkaContainer
	wg.Go(func() {
		kafkaContainer = must(kafka.Run(ctx,
			"confluentinc/confluent-local:7.5.0",
			kafka.WithClusterID("test-cluster"),
			testcontainers.WithEnv(map[string]string{
				"KAFKA_AUTO_CREATE_TOPICS_ENABLE": "true",
			})))
	})
	wg.Wait()

	brokers := must(kafkaContainer.Brokers(ctx))
	kafkaCluster = KafkaTestCluster{
		container: kafkaContainer,
		brokers:   brokers,
	}

	code := m.Run()

	go func() { _ = testcontainers.TerminateContainer(cassandraDB.container) }()
	go func() { _ = testcontainers.TerminateContainer(kafkaCluster.container) }()

	os.Exit(code)
}

func must[T any](v T, err error) T {
	if err != nil {
		log.Fatalf("setup failed: %v", err)
	}
	return v
}
