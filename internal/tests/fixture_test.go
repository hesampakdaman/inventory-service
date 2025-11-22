package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
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

	server *httptest.Server
	client *http.Client
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
		server:  server,
		client:  server.Client(),
	}
}

func (f Fixture) DoRequest(
	method, path string,
	body io.Reader,
	headers map[string]string,
) (*http.Response, []byte, error) {
	if f.server == nil || f.client == nil {
		return nil, nil, fmt.Errorf("fixture HTTP server not initialized")
	}

	req, err := http.NewRequest(method, f.server.URL+path, body)
	if err != nil {
		return nil, nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := f.client.Do(req)
	if err != nil {
		return resp, nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp, nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return resp, data, nil
}

func (f Fixture) DoJSONRaw(method, path string, body any) (*http.Response, []byte, error) {
	var reqBody io.Reader
	if body != nil {
		var buf bytes.Buffer
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			return nil, nil, fmt.Errorf("failed to encode request body: %w", err)
		}
		reqBody = &buf
	}

	headers := map[string]string{
		"Accept": "application/json",
	}
	if body != nil {
		headers["Content-Type"] = "application/json"
	}

	return f.DoRequest(method, path, reqBody, headers)
}

func (f Fixture) DoJSON(method, path string, body any, v any) (*http.Response, error) {
	resp, data, err := f.DoJSONRaw(method, path, body)
	if err != nil {
		return resp, err
	}

	if v == nil {
		return resp, nil
	}

	if err := json.Unmarshal(data, v); err != nil {
		return resp, fmt.Errorf(
			"failed to decode JSON from %s %s: %w\nresponse body: %s",
			method,
			path,
			err,
			data,
		)
	}

	return resp, nil
}
