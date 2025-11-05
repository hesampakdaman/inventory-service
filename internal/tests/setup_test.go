package tests

import (
	"context"
	"fmt"
	"log"
	"testing"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
	"github.com/docker/go-connections/nat"
	"github.com/google/uuid"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/cassandra"
)

var cassandraDB CassandraTestDB

type CassandraTestDB struct {
	cluster *gocql.ClusterConfig
	host    string
	port    nat.Port
}

func (c CassandraTestDB) Session(t *testing.T) (*gocql.Session, error) {
	ctx := t.Context()

	cluster := *c.cluster
	cluster.Keyspace = fmt.Sprintf("test_%s", uuid.NewString())
	session, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}

	session.Query(fmt.Sprintf(`
        CREATE KEYSPACE IF NOT EXISTS %s
        WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1}`,
		cluster.Keyspace)).
		ExecContext(ctx)
	return session, nil
}

func TestMain(m *testing.M) {
	ctx := context.Background()

	cassandraContainer := must(cassandra.Run(ctx, "cassandra:4.1.3"))
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
