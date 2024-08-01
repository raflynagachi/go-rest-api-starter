package encoder

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeJson(t *testing.T) {
	data := map[string]string{"key": "value"}

	recorder := httptest.NewRecorder()

	err := EncodeJson(recorder, data)
	if err != nil {
		t.Fatalf("EncodeJson() error = %v", err)
	}

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("EncodeJson() status = %v, want %v", status, http.StatusOK)
	}

	if contentType := recorder.Header().Get("Content-Type"); contentType != "application/json" {
		t.Errorf("EncodeJson() Content-Type = %v, want %v", contentType, "application/json")
	}

	expectedBody, _ := json.Marshal(data)
	body := recorder.Body.String()
	body = strings.ReplaceAll(body, "\n", "")
	if body != string(expectedBody) {
		t.Errorf("EncodeJson() body = %v, want %v", body, string(expectedBody))
	}
}

func TestDecodeJson(t *testing.T) {
	data := map[string]string{"key": "value"}
	dataBytes, _ := json.Marshal(data)

	req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader(dataBytes))
	req.Header.Set("Content-Type", "application/json")

	var result map[string]string

	err := DecodeJson(req, &result)
	if err != nil {
		t.Fatalf("DecodeJson() error = %v", err)
	}

	if !assert.Equal(t, data, result) {
		t.Errorf("DecodeJson() result = %v, want %v", result, data)
	}
}
