package tests

import (
	"fmt"
	"io"
	"log/slog"
	"net/http/httptest"
	"os"
	"testing"

	gocql "github.com/apache/cassandra-gocql-driver/v2"

	"github.com/hesampakdaman/inventory-service/internal/adapters/kafka"
	"github.com/hesampakdaman/inventory-service/internal/app"
	"github.com/hesampakdaman/inventory-service/internal/core/messagebus"
)

type Fixture struct {
	Session *gocql.Session
	App     *app.App
	Bus     *messagebus.Bus
}

func NewFixture(t *testing.T) Fixture {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	session, sessErr := cassandraDB.Session(t)
	if sessErr != nil {
		logs, err := cassandraDB.container.Logs(t.Context())
		if err != nil {
			t.Fatalf("failed to get logs: %v", err)
		}
		defer logs.Close()
		buf, _ := io.ReadAll(logs)
		fmt.Println(string(buf))
		t.Fatalf("%s", sessErr)
	}
	t.Cleanup(func() { session.Close() })

	kClient, topic := kafkaCluster.NewTopicClient(t)
	t.Cleanup(func() { kClient.Close() })

	application := app.New(logger, session, kClient, kafka.Topic(topic))
	go func() { application.Consumer.Consume(t.Context()) }()

	server := httptest.NewServer(application.Router)
	t.Cleanup(server.Close)

	return Fixture{
		Session: session,
		App:     application,
		Bus:     application.Bus,
	}
}
