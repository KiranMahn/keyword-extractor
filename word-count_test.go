package main

import (
	"regexp"
	"testing"
)

/*
This file tests for:
- correct word count is returned
- repeated words are counted correctly
- case insensitivity of word counting
- small words are not counted
- filtering out words starting with digits
- empty content handling
- handling content with only stopwords and short words
- punctuation handling
- handling content without stopwords present
- it works with custom word splitters
- benchmark for getWordCount function
*/
func TestGetWordCount(t *testing.T) {
	// Setup common test data
	stopwords := map[string]struct{}{
		"the": {},
		"and": {},
		"of":  {},
		"to":  {},
		"a":   {},
		"in":  {},
		"is":  {},
		"it":  {},
	}

	// Common word splitter regex for splitting on non-word characters
	wordSplitter := regexp.MustCompile(`\W+`)

	// Test that the correct word count is returned
	t.Run("BasicWordCounting", func(t *testing.T) {
		content := "The quick brown fox jumps over the lazy dog"
		result, err := getWordCount(content, stopwords, wordSplitter)

		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		// Expected words (excluding stopwords "the" and short words)
		expectedWords := map[string]int{
			"quick": 1,
			"brown": 1,
			"fox":   1,
			"jumps": 1,
			"over":  1,
			"lazy":  1,
			"dog":   1,
		}

		if len(result) != len(expectedWords) {
			t.Errorf("Expected %d words, got %d", len(expectedWords), len(result))
		}

		// Verify each expected word is present with the correct count
		for word, expectedCount := range expectedWords {
			if actualCount, exists := result[word]; !exists {
				t.Errorf("Expected word '%s' to be present", word)
			} else if actualCount != expectedCount {
				t.Errorf("Expected count %d for word '%s', got %d", expectedCount, word, actualCount)
			}
		}

		// Verify stopwords are not included
		for stopword := range stopwords {
			if _, exists := result[stopword]; exists {
				t.Errorf("Stopword '%s' should not be included", stopword)
			}
		}
	})

	// Test handling of repeated words
	t.Run("RepeatedWords", func(t *testing.T) {
		content := "apple banana apple cherry banana apple"
		result, err := getWordCount(content, stopwords, wordSplitter)

		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		expected := map[string]int{
			"apple":  3,
			"banana": 2,
			"cherry": 1,
		}

		// Check if the result matches the expected counts
		for word, expectedCount := range expected {
			if actualCount, exists := result[word]; !exists {
				t.Errorf("Expected word '%s' to be present", word)
			} else if actualCount != expectedCount {
				t.Errorf("Expected count %d for word '%s', got %d", expectedCount, word, actualCount)
			}
		}
	})

	// Test case insensitivity
	t.Run("CaseInsensitive", func(t *testing.T) {
		content := "Apple APPLE apple ApPlE"
		result, err := getWordCount(content, stopwords, wordSplitter)

		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		if count, exists := result["apple"]; !exists {
			t.Error("Expected 'apple' to be present")
		} else if count != 4 {
			t.Errorf("Expected count 4 for 'apple', got %d", count)
		}

		// Verify no uppercase versions exist
		upperCaseVariants := []string{"Apple", "APPLE", "ApPlE"}
		for _, variant := range upperCaseVariants {
			if _, exists := result[variant]; exists {
				t.Errorf("Uppercase variant '%s' should not exist (should be converted to lowercase)", variant)
			}
		}
	})

	// Test that small words are not counted
	t.Run("FilterShortWords", func(t *testing.T) {
		content := "a an the programming go is fun"
		result, err := getWordCount(content, stopwords, wordSplitter)

		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		// Words with length <= 2 should be filtered out
		shortWords := []string{"a", "an", "go", "is"}
		for _, word := range shortWords {
			if _, exists := result[word]; exists {
				t.Errorf("Short word '%s' should be filtered out", word)
			}
		}

		// Long words should be included (excluding stopwords)
		if _, exists := result["programming"]; !exists {
			t.Error("Expected 'programming' to be present")
		}
		if _, exists := result["fun"]; !exists {
			t.Error("Expected 'fun' to be present")
		}
	})

	// Test for filtering out words starting with digits
	t.Run("FilterDigitWords", func(t *testing.T) {
		content := "123 456 abc 789def hello 2023 world"
		result, err := getWordCount(content, stopwords, wordSplitter)

		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		// Words starting with digits should be filtered out
		digitWords := []string{"123", "456", "789def", "2023"}
		for _, word := range digitWords {
			if _, exists := result[word]; exists {
				t.Errorf("Word starting with digit '%s' should be filtered out", word)
			}
		}

		// Non-digit words should be included
		expectedWords := []string{"abc", "hello", "world"}
		for _, word := range expectedWords {
			if _, exists := result[word]; !exists {
				t.Errorf("Expected word '%s' to be present", word)
			}
		}
	})

	// Test for empty content
	t.Run("EmptyContent", func(t *testing.T) {
		content := ""
		_, err := getWordCount(content, stopwords, wordSplitter)
		if err == nil {
			t.Fatalf("Expected error, got: %v", err)
		}
	})

	// Test for content with only stopwords and short words
	t.Run("OnlyStopwordsAndShortWords", func(t *testing.T) {
		content := "the and of to a in is it an I"
		_, err := getWordCount(content, stopwords, wordSplitter)

		if err == nil {
			t.Fatalf("Expected error, got: %v", err)
		}

	})

	// Test for punctuation handling
	t.Run("PunctuationHandling", func(t *testing.T) {
		content := "Hello, world! How are you? I'm fine."
		result, err := getWordCount(content, stopwords, wordSplitter)
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		expectedWords := []string{"hello", "world", "how", "are", "you", "fine"}
		for _, word := range expectedWords {
			if _, exists := result[word]; !exists {
				t.Errorf("Expected word '%s' to be present (punctuation should be removed)", word)
			}
		}

		// Verify punctuation doesn't appear as separate words
		punctuation := []string{",", "!", "?", ".", "'m"}
		for _, punct := range punctuation {
			if _, exists := result[punct]; exists {
				t.Errorf("Punctuation '%s' should not appear as a word", punct)
			}
		}
	})

	// Test for handling content without stopwords present
	t.Run("EmptyStopwords", func(t *testing.T) {
		content := "the quick brown fox"
		emptyStopwords := make(map[string]struct{})
		result, err := getWordCount(content, emptyStopwords, wordSplitter)
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		// With no stopwords, all words longer than 2 chars should be included
		expectedWords := []string{"the", "quick", "brown", "fox"}
		for _, word := range expectedWords {
			if _, exists := result[word]; !exists {
				t.Errorf("Expected word '%s' to be present when no stopwords", word)
			}
		}
	})

	// Test it works with custom word splitters
	t.Run("CustomWordSplitter", func(t *testing.T) {
		content := "word1-word2_word3 word4"
		// Custom splitter that splits on hyphens and underscores too
		customSplitter := regexp.MustCompile(`[\s\-_]+`)
		result, err := getWordCount(content, stopwords, customSplitter)
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		expectedWords := []string{"word1", "word2", "word3", "word4"}
		for _, word := range expectedWords {
			if _, exists := result[word]; !exists {
				t.Errorf("Expected word '%s' to be present with custom splitter", word)
			}
		}
	})
}

// Benchmark test for getWordCount function
func BenchmarkGetWordCount(b *testing.B) {
	content := "The quick brown fox jumps over the lazy dog. " +
		"This is a sample text for benchmarking the word count function. " +
		"It contains various words that should be counted properly."

	stopwords := map[string]struct{}{
		"the": {}, "and": {}, "of": {}, "to": {}, "a": {}, "in": {}, "is": {}, "it": {},
	}

	wordSplitter := regexp.MustCompile(`\W+`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		getWordCount(content, stopwords, wordSplitter)
	}
}
