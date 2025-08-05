package main

import (
	"os"
	"path/filepath"
	"testing"
)

/*
This file tests for:
- loading a vald stopwords file
- loading a non-existent file
- loading an empty file
- loading a file with whitespace and varied formatting
- loading a file with only a single word
- benchmark for loading stopwords
*/
func TestLoadStopwords(t *testing.T) {
	// Test loading a valid stopwords file
	t.Run("ValidStopwordsFile", func(t *testing.T) {
		// check no error from loading valid stopwords file
		stopwords, err := LoadStopwords("data/stopwords.txt")
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		// check valid stopwords file is not empty
		if len(stopwords) == 0 {
			t.Fatal("Expected stopwords to be loaded, got empty map")
		}

		// Check if some common stopwords are present
		expectedWords := []string{"the", "of", "to", "and", "a"}
		for _, word := range expectedWords {
			if _, exists := stopwords[word]; !exists {
				t.Errorf("Expected stopword '%s' to be present", word)
			}
		}
	})

	// Test for a non-existent file
	t.Run("NonExistentFile", func(t *testing.T) {
		// load nonexistent file
		stopwords, err := LoadStopwords("non_existent_file.txt")
		// check if error is returned
		if err == nil {
			t.Fatal("Expected error for non-existent file, got nil")
		}
		// check that no stopwords are returned
		if stopwords != nil {
			t.Error("Expected nil stopwords map for non-existent file")
		}
	})

	// Test for an Empty file
	t.Run("EmptyFile", func(t *testing.T) {
		// Create a temporary empty file
		tempDir := t.TempDir()
		emptyFile := filepath.Join(tempDir, "empty.txt")
		file, err := os.Create(emptyFile)
		if err != nil {
			t.Fatalf("Failed to create temporary file: %v", err)
		}
		file.Close()

		// Load stopwords from the empty file
		stopwords, err := LoadStopwords(emptyFile)

		// Check if error is returned
		if err == nil {
			t.Fatalf("Expected error for empty file, got: %v", err)
		}

		// Check if stopwords map is empty
		if len(stopwords) != 0 {
			t.Errorf("Expected empty map for empty file, got %d entries", len(stopwords))
		}
	})

	// Test function can handle a file with whitespace and varied formatting
	t.Run("FileWithWhitespace", func(t *testing.T) {
		// Create a temporary file with various whitespace scenarios
		tempDir := t.TempDir()
		testFile := filepath.Join(tempDir, "test_stopwords.txt")
		content := "  word1  \n\tword2\t\n\nword3\n  \n\tword4\t  \n"

		// Create the test file
		err := os.WriteFile(testFile, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		// attempt to load stopwords
		stopwords, err := LoadStopwords(testFile)
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		// Check if stopwords are loaded correctly
		expectedWords := []string{"word1", "word2", "word3", "word4"}
		if len(stopwords) != len(expectedWords) {
			t.Errorf("Expected %d words, got %d", len(expectedWords), len(stopwords))
		}

		for _, word := range expectedWords {
			if _, exists := stopwords[word]; !exists {
				t.Errorf("Expected word '%s' to be present", word)
			}
		}

		// Ensure empty lines don't create entries
		if _, exists := stopwords[""]; exists {
			t.Error("Empty string should not be in stopwords map")
		}
	})

	// Test for a file with only a single word
	t.Run("SingleWordFile", func(t *testing.T) {
		// create a file with a single word
		tempDir := t.TempDir()
		singleWordFile := filepath.Join(tempDir, "single.txt")

		err := os.WriteFile(singleWordFile, []byte("hello"), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		// Load stopwords from the single word file
		stopwords, err := LoadStopwords(singleWordFile)
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		// Check if stopwords map contains the single word and that word is hello
		if len(stopwords) != 1 {
			t.Errorf("Expected 1 word, got %d", len(stopwords))
		}

		if _, exists := stopwords["hello"]; !exists {
			t.Error("Expected 'hello' to be present")
		}
	})
}

// Benchmark test for LoadStopwords function
func BenchmarkLoadStopwords(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := LoadStopwords("data/stopwords.txt")
		if err != nil {
			b.Fatalf("Benchmark failed: %v", err)
		}
	}
}
