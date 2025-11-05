package tests

import (
	"fmt"
	"testing"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
)

func TestLeTest(t *testing.T) {
	fx := NewFixture(t)
	batch := fx.Session.Batch(gocql.LoggedBatch)

	batch.Query(`
    INSERT INTO products (product_id, available, reserved, title, description)
    VALUES (?, ?, ?, ?, ?);`,
		"1", 2, 3, "Title", "Desc",
	)

	if err := batch.ExecContext(t.Context()); err != nil {
		fmt.Println(err)
		t.Fatalf("%s", err)
	}
}
