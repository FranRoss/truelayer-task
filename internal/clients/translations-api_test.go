//go:build unit

package clients

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func setupTranslationMockServer(handler http.HandlerFunc) (*httptest.Server, *TranslationClient) {
	server := httptest.NewServer(handler)
	client := &TranslationClient{
		baseURL:    server.URL,
		httpClient: server.Client(),
	}
	return server, client
}

func TestTranslate_Success(t *testing.T) {
	mockResponseJSON := `{
		"success": {"total": 1},
		"contents": {
			"translated": "Lost a planet, Master Obi-Wan has.",
			"text": "Master Obi-Wan has lost a planet.",
			"translation": "yoda"
		}
	}`

	server, client := setupTranslationMockServer(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockResponseJSON))
	})
	defer server.Close()

	result, err := client.Translate("Master Obi-Wan has lost a planet.", Yoda)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	expectedTranslation := "Lost a planet, Master Obi-Wan has."
	if result != expectedTranslation {
		t.Errorf("expected translation %q, got %q", expectedTranslation, result)
	}
}

func TestTranslate_NetworkFailureError(t *testing.T) {
	server, client := setupTranslationMockServer(func(w http.ResponseWriter, r *http.Request) {
		hj, ok := w.(http.Hijacker)
		if !ok {
			t.Fatal("webserver doesn't support hijacking")
		}
		conn, _, _ := hj.Hijack()
		conn.Close()
	})
	defer server.Close()

	_, err := client.Translate("Hello world", Yoda)
	if err == nil {
		t.Fatal("expected a network failure error, but got nil")
	}

	if !strings.Contains(err.Error(), "failed making request to FunTranslations") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestTranslate_Non200StatusCode(t *testing.T) {
	server, client := setupTranslationMockServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
	})
	defer server.Close()

	_, err := client.Translate("Hello world", Yoda)
	if err == nil {
		t.Fatal("expected an HTTP status error, but got nil")
	}

	expectedErr := fmt.Sprintf("http error status %d", http.StatusTooManyRequests)
	if err.Error() != expectedErr {
		t.Errorf("expected error %q, got %q", expectedErr, err.Error())
	}
}

func TestTranslate_APIErrorDetailsPayload(t *testing.T) {
	mockResponseJSON := `{
		"error": {
			"code": 429,
			"message": "Too Many Requests: Rate limit exceeded"
		}
	}`

	server, client := setupTranslationMockServer(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK) // Even though it's an API error, client expects 200 OK before parsing it
		w.Write([]byte(mockResponseJSON))
	})
	defer server.Close()

	_, err := client.Translate("Hello world", Yoda)
	if err == nil {
		t.Fatal("expected an API details payload error, but got nil")
	}

	expectedErr := "429 Too Many Requests: Rate limit exceeded"
	if err.Error() != expectedErr {
		t.Errorf("expected error %q, got %q", expectedErr, err.Error())
	}
}
