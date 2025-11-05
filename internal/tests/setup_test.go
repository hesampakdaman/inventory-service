package tests

import (
	"context"
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/apache/cassandra-gocql-driver/v2"
	"github.com/docker/go-connections/nat"
	"github.com/google/uuid"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/cassandra"
	"github.com/testcontainers/testcontainers-go/wait"
)

var cassandraDB CassandraTestDB
var cassandraContainer *cassandra.CassandraContainer

type CassandraTestDB struct {
	cluster *gocql.ClusterConfig
	host    string
	port    nat.Port
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

	return session, nil
}

func TestMain(m *testing.M) {
	ctx := context.Background()

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
		host:    host,
		port:    port,
	}

	m.Run()

	_ = testcontainers.TerminateContainer(cassandraContainer)
}

func must[T any](v T, err error) T {
	if err != nil {
		log.Fatalf("setup failed: %v", err)
	}
	return v
}
