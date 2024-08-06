package router

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/raflynagachi/go-rest-api-starter/config"
	"github.com/raflynagachi/go-rest-api-starter/internal/handler/definition/mocks"
	"github.com/raflynagachi/go-rest-api-starter/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	mockHandler = new(mocks.APIHandler)
	mockCfg     = &config.Config{
		App: config.App{
			Port: 8080,
		},
	}
	mockLogger = logger.NewLogger()
)

func TestNewRouter(t *testing.T) {
	router := New(mockCfg, mockLogger, mockHandler)

	assert.NotNil(t, router)
	assert.Equal(t, mockCfg, router.Cfg)
	assert.NotNil(t, router.Router)
}

func TestRouter_ServeHTTP(t *testing.T) {
	r := New(mockCfg, mockLogger, mockHandler)

	go func() {
		r.ServeHTTP()
	}()

	// allow the server to start
	time.Sleep(100 * time.Millisecond)

	resp, err := http.Get("http://localhost:8080/ping")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// test shutdown
	err = r.Shutdown(context.Background())
	require.NoError(t, err)
}

func TestRouter_Shutdown(t *testing.T) {
	r := New(mockCfg, mockLogger, mockHandler)

	err := r.Shutdown(context.Background())
	require.NoError(t, err)

	resp, err := http.Get("http://localhost:8080/ping")
	assert.Error(t, err)
	assert.Nil(t, resp)
}
