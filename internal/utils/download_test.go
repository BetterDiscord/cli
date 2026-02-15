package utils

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestDownloadFile(t *testing.T) {
	// Create a test server
	testData := []byte("test file content")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request headers
		if r.Header.Get("User-Agent") != "BetterDiscord/cli" {
			t.Errorf("Expected User-Agent header 'BetterDiscord/cli', got '%s'", r.Header.Get("User-Agent"))
		}
		if r.Header.Get("Accept") != "application/octet-stream" {
			t.Errorf("Expected Accept header 'application/octet-stream', got '%s'", r.Header.Get("Accept"))
		}

		w.WriteHeader(http.StatusOK)
		w.Write(testData)
	}))
	defer server.Close()

	// Create temp directory for test file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "downloaded.txt")

	// Download the file
	resp, err := DownloadFile(server.URL, testFile)
	if err != nil {
		t.Fatalf("DownloadFile() failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	// Verify file was created
	if !Exists(testFile) {
		t.Errorf("Downloaded file does not exist: %s", testFile)
	}

	// Verify file content
	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read downloaded file: %v", err)
	}

	if string(content) != string(testData) {
		t.Errorf("File content mismatch. Expected '%s', got '%s'", string(testData), string(content))
	}
}

func TestDownloadFile_BadStatusCode(t *testing.T) {
	// Create a test server that returns 404
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "downloaded.txt")

	_, err := DownloadFile(server.URL, testFile)
	if err == nil {
		t.Error("DownloadFile() should have returned an error for 404 status")
	}
}

func TestDownloadFile_InvalidURL(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "downloaded.txt")

	_, err := DownloadFile("http://invalid.test.nonexistent.domain", testFile)
	if err == nil {
		t.Error("DownloadFile() should have returned an error for invalid URL")
	}
}

func TestDownloadFile_InvalidPath(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test"))
	}))
	defer server.Close()

	// Try to write to an invalid path (directory doesn't exist)
	invalidPath := "/nonexistent/directory/file.txt"
	_, err := DownloadFile(server.URL, invalidPath)
	if err == nil {
		t.Error("DownloadFile() should have returned an error for invalid file path")
	}
}

func TestDownloadJSON(t *testing.T) {
	// Define a test struct
	type TestData struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	expectedData := TestData{
		Name:  "test",
		Value: 42,
	}

	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request headers
		if r.Header.Get("User-Agent") != "BetterDiscord/cli" {
			t.Errorf("Expected User-Agent header 'BetterDiscord/cli', got '%s'", r.Header.Get("User-Agent"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(expectedData)
	}))
	defer server.Close()

	// Download and parse JSON
	result, err := DownloadJSON[TestData](server.URL)
	if err != nil {
		t.Fatalf("DownloadJSON() failed: %v", err)
	}

	if result.Name != expectedData.Name {
		t.Errorf("Expected Name '%s', got '%s'", expectedData.Name, result.Name)
	}

	if result.Value != expectedData.Value {
		t.Errorf("Expected Value %d, got %d", expectedData.Value, result.Value)
	}
}

func TestDownloadJSON_BadStatusCode(t *testing.T) {
	type TestData struct {
		Name string `json:"name"`
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	_, err := DownloadJSON[TestData](server.URL)
	if err == nil {
		t.Error("DownloadJSON() should have returned an error for 500 status")
	}
}

func TestDownloadJSON_InvalidJSON(t *testing.T) {
	type TestData struct {
		Name string `json:"name"`
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid json {"))
	}))
	defer server.Close()

	_, err := DownloadJSON[TestData](server.URL)
	if err == nil {
		t.Error("DownloadJSON() should have returned an error for invalid JSON")
	}
}

func TestDownloadJSON_ComplexStruct(t *testing.T) {
	type NestedData struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	type ComplexData struct {
		Items     []NestedData  `json:"items"`
		Timestamp time.Time     `json:"timestamp"`
		Active    bool          `json:"active"`
	}

	now := time.Now().UTC().Round(time.Second)
	expectedData := ComplexData{
		Items: []NestedData{
			{ID: 1, Name: "first"},
			{ID: 2, Name: "second"},
		},
		Timestamp: now,
		Active:    true,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(expectedData)
	}))
	defer server.Close()

	result, err := DownloadJSON[ComplexData](server.URL)
	if err != nil {
		t.Fatalf("DownloadJSON() failed: %v", err)
	}

	if len(result.Items) != len(expectedData.Items) {
		t.Errorf("Expected %d items, got %d", len(expectedData.Items), len(result.Items))
	}

	if result.Active != expectedData.Active {
		t.Errorf("Expected Active %v, got %v", expectedData.Active, result.Active)
	}

	if !result.Timestamp.Equal(expectedData.Timestamp) {
		t.Errorf("Expected Timestamp %v, got %v", expectedData.Timestamp, result.Timestamp)
	}
}

func TestDownloadJSON_InvalidURL(t *testing.T) {
	type TestData struct {
		Name string `json:"name"`
	}

	_, err := DownloadJSON[TestData]("http://invalid.test.nonexistent.domain")
	if err == nil {
		t.Error("DownloadJSON() should have returned an error for invalid URL")
	}
}
