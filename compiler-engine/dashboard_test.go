package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDashboardStatus(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/status", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleStatus)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, status)
	}

	var resp StatusResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if resp.Version == "" {
		t.Errorf("Expected non-empty version in status response")
	}
	if len(resp.RuntimeHeaders) == 0 {
		t.Errorf("Expected runtime headers list in status response")
	}
}

func TestDashboardTranspile(t *testing.T) {
	javaSrc := `
public class PingApp {
    public int ping() {
        return 42;
    }
}
`
	reqBody, _ := json.Marshal(TranspileRequest{Source: javaSrc})
	req, err := http.NewRequest("POST", "/api/transpile", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleTranspileSnippet)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, status)
	}

	var resp TranspileResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to decode transpile response: %v", err)
	}

	if !resp.Success {
		t.Fatalf("Expected transpile success, got error: %s", resp.Error)
	}
	if !strings.Contains(resp.Header, "class PingApp : public java::lang::Object") {
		t.Errorf("Expected PingApp class declaration in header, got:\n%s", resp.Header)
	}
	if !strings.Contains(resp.Cpp, "int32_t PingApp::ping()") {
		t.Errorf("Expected ping method definition in cpp, got:\n%s", resp.Cpp)
	}
}
