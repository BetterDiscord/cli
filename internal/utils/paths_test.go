package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExists(t *testing.T) {
	// Create a temporary file for testing
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "testfile.txt")

	// Test non-existent file
	if Exists(tmpFile) {
		t.Errorf("Exists() returned true for non-existent file: %s", tmpFile)
	}

	// Create the file
	if err := os.WriteFile(tmpFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test existing file
	if !Exists(tmpFile) {
		t.Errorf("Exists() returned false for existing file: %s", tmpFile)
	}

	// Test directory
	if !Exists(tmpDir) {
		t.Errorf("Exists() returned false for existing directory: %s", tmpDir)
	}

	// Test non-existent directory
	nonExistentDir := filepath.Join(tmpDir, "nonexistent", "directory")
	if Exists(nonExistentDir) {
		t.Errorf("Exists() returned true for non-existent directory: %s", nonExistentDir)
	}
}

func TestFilter(t *testing.T) {
	tests := []struct {
		name       string
		input      []int
		filterFunc func(int) bool
		expected   []int
	}{
		{
			name:       "Filter even numbers",
			input:      []int{1, 2, 3, 4, 5, 6},
			filterFunc: func(n int) bool { return n%2 == 0 },
			expected:   []int{2, 4, 6},
		},
		{
			name:       "Filter numbers greater than 5",
			input:      []int{1, 3, 5, 7, 9, 11},
			filterFunc: func(n int) bool { return n > 5 },
			expected:   []int{7, 9, 11},
		},
		{
			name:       "Filter all items (return empty)",
			input:      []int{1, 2, 3},
			filterFunc: func(n int) bool { return false },
			expected:   []int{},
		},
		{
			name:       "Filter none (return all)",
			input:      []int{1, 2, 3},
			filterFunc: func(n int) bool { return true },
			expected:   []int{1, 2, 3},
		},
		{
			name:       "Empty input",
			input:      []int{},
			filterFunc: func(n int) bool { return true },
			expected:   []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Filter(tt.input, tt.filterFunc)

			if len(result) != len(tt.expected) {
				t.Errorf("Filter() returned slice of length %d, expected %d", len(result), len(tt.expected))
				return
			}

			for i, v := range result {
				if v != tt.expected[i] {
					t.Errorf("Filter() result[%d] = %v, expected %v", i, v, tt.expected[i])
				}
			}
		})
	}
}

func TestFilterWithStrings(t *testing.T) {
	input := []string{"apple", "banana", "cherry", "date"}
	filterFunc := func(s string) bool {
		return len(s) > 5
	}
	expected := []string{"banana", "cherry"}

	result := Filter(input, filterFunc)

	if len(result) != len(expected) {
		t.Errorf("Filter() returned slice of length %d, expected %d", len(result), len(expected))
		return
	}

	for i, v := range result {
		if v != expected[i] {
			t.Errorf("Filter() result[%d] = %v, expected %v", i, v, expected[i])
		}
	}
}

func TestFilterWithStructs(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	input := []Person{
		{Name: "Alice", Age: 25},
		{Name: "Bob", Age: 17},
		{Name: "Charlie", Age: 30},
		{Name: "Diana", Age: 16},
	}

	// Filter adults (age >= 18)
	filterFunc := func(p Person) bool {
		return p.Age >= 18
	}

	result := Filter(input, filterFunc)

	if len(result) != 2 {
		t.Errorf("Filter() returned %d persons, expected 2 adults", len(result))
		return
	}

	if result[0].Name != "Alice" || result[1].Name != "Charlie" {
		t.Errorf("Filter() returned unexpected persons: %v", result)
	}
}
