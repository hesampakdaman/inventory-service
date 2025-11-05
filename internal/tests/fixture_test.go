package tests

import (
	"testing"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
)

type Fixture struct {
	Session *gocql.Session
}

func NewFixture(t *testing.T) Fixture {
	session, err := cassandraDB.Session(t)
	if err != nil {
		t.Fatal()
	}
	defer t.Cleanup(func() { session.Close() })

	return Fixture{
		Session: session,
	}
}
