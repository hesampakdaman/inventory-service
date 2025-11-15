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
	cassandraDB        CassandraTestDB
	cassandraContainer *cassandra.CassandraContainer
	kafkaContainer     *kafka.KafkaContainer
)

type CassandraTestDB struct {
	cluster *gocql.ClusterConfig
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
		cassandraContainer = must(
			cassandra.Run(
				ctx,
				"cassandra:4.1.3",
				testcontainers.WithWaitStrategy(wait.ForLog("Starting listening for CQL clients")),
			),
		)
		host := must(cassandraContainer.Host(ctx))
		port := must(cassandraContainer.MappedPort(ctx, "9042/tcp"))

		cassandraDB = CassandraTestDB{
			cluster: gocql.NewCluster(fmt.Sprintf("%s:%s", host, port.Port())),
		}
	})

	wg.Go(func() {
		kafkaContainer = must(kafka.Run(ctx,
			"confluentinc/confluent-local:7.5.0",
			kafka.WithClusterID("test-cluster"),
		))
	})
	wg.Wait()

	brokers, err := kafkaContainer.Brokers(ctx)
	if err != nil {
		log.Printf("failed to start container: %s", err)
		return
	}

	cl, err := kgo.NewClient(
		kgo.SeedBrokers(brokers...),
		kgo.ConsumerGroup("inventory-service"),
		kgo.ConsumeTopics("inventory"),
	)

	fmt.Printf("%v", cl)

	m.Run()

	go func() { _ = testcontainers.TerminateContainer(cassandraContainer) }()
	go func() { _ = testcontainers.TerminateContainer(kafkaContainer) }()
}

func must[T any](v T, err error) T {
	if err != nil {
		log.Fatalf("setup failed: %v", err)
	}
	return v
}
