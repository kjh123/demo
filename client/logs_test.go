package main

import (
	"context"
	"data/logs"
	"encoding/json"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/labstack/gommon/random"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tdewolff/parse/v2/buffer"
)

type mockLogWriter struct {
	mock.Mock
}

func (m *mockLogWriter) Writer(ctx context.Context, log logs.BehaviorLog) error {
	args := m.Called(ctx, log)
	return args.Error(0)
}

func TestLogService_Writer(t *testing.T) {
	writer := &mockLogWriter{}
	writer.On("Writer", mock.Anything, mock.Anything).Return(nil)

	service := NewLogService(writer)
	w := httptest.NewRecorder()

	log := logs.BehaviorLog{
		UID:  rand.Int63n(1000000),
		UA:   random.String(10),
		IP:   random.String(10),
		Tags: []string{"tag1", "tag2", "tag3", "tag4", "tag5"},
	}

	b, err := json.Marshal(log)
	assert.NoError(t, err)
	req := httptest.NewRequest("POST", "/_/logs", buffer.NewReader(b))
	app := gin.New()
	service.Mount(app.Group("_"))
	app.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, w.Body.String())
}

func TestLogService_Write(t *testing.T) {
	host := os.Getenv("INFLUX_HOST")
	token := os.Getenv("INFLUX_TOKEN")
	if host == "" || token == "" {
		t.Skip("INFLUX_HOST or INFLUX_TOKEN is not set")
	}

	client := influxdb2.NewClient(host, token)
	ok, err := client.Ping(context.Background())
	assert.True(t, ok)
	assert.NoError(t, err)

	writer := logs.NewInfluxRepository(&logs.InfluxClient{
		Client:   client,
		WriteAPI: client.WriteAPIBlocking("org", "bucket"),
	})
	service := NewLogService(writer)
	w := httptest.NewRecorder()

	log := logs.BehaviorLog{
		UID:  rand.Int63n(1000000),
		UA:   random.String(10),
		IP:   random.String(10),
		Tags: []string{"tag1", "tag2", "tag3", "tag4", "tag5"},
	}

	b, err := json.Marshal(log)
	assert.NoError(t, err)
	req := httptest.NewRequest("POST", "/_/logs", buffer.NewReader(b))
	app := gin.New()
	service.Mount(app.Group("_"))
	app.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, w.Body.String())
}
