package tests

import (
	"fmt"
	"io"
	"log/slog"
	"net/http/httptest"
	"os"
	"testing"

	gocql "github.com/apache/cassandra-gocql-driver/v2"

	"github.com/hesampakdaman/inventory-service/internal/app"
)

type Fixture struct {
	Session *gocql.Session
	App     *app.App
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

	application := app.New(logger, session, kafkaClient)

	server := httptest.NewServer(application.Router)
	t.Cleanup(server.Close)

	return Fixture{
		Session: session,
		App:     application,
	}
}
