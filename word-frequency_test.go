package main

import (
	"math"
	"regexp"
	"testing"
)

/*
This file tests for:
- correct word frequency calculation
- content with a single word
- empty content handling
- handling a non-existent word in TCI
- handling punctuation in content
- handling case sensitivity

*/

func TestGetWordFrequency(t *testing.T) {
	// Common word splitter regex for splitting on non-word characters
	wordSplitter := regexp.MustCompile(`\W+`)

	t.Run("BasicWordFrequency", func(t *testing.T) {
		content := "apple banana apple cherry"

		// Create term count index
		tci := TermCountIndex{
			"apple":  2,
			"banana": 1,
			"cherry": 1,
		}

		result := getWordFrequency(content, wordSplitter, tci)

		// Total words = 4, so frequencies should be:
		// apple: 2/4 = 0.5
		// banana: 1/4 = 0.25
		// cherry: 1/4 = 0.25
		expectedFrequencies := map[string]float64{
			"apple":  0.5,
			"banana": 0.25,
			"cherry": 0.25,
		}

		for word, expectedFreq := range expectedFrequencies {
			if actualFreq, exists := result[word]; !exists {
				t.Errorf("Expected word '%s' to be present", word)
			} else if math.Abs(actualFreq-expectedFreq) > 0.0001 {
				t.Errorf("Expected frequency %.4f for word '%s', got %.4f", expectedFreq, word, actualFreq)
			}
		}
	})

	// Test for handling for content with a single word
	t.Run("SingleWord", func(t *testing.T) {
		content := "hello"

		tci := TermCountIndex{
			"hello": 1,
		}

		result := getWordFrequency(content, wordSplitter, tci)

		// Only one word, so frequency should be 1.0
		if freq, exists := result["hello"]; !exists {
			t.Error("Expected 'hello' to be present")
		} else if math.Abs(freq-1.0) > 0.0001 {
			t.Errorf("Expected frequency 1.0 for single word, got %.4f", freq)
		}
	})

	// Test for handling empty content
	t.Run("EmptyContent", func(t *testing.T) {
		content := ""
		tci := TermCountIndex{}

		result := getWordFrequency(content, wordSplitter, tci)

		if len(result) != 0 {
			t.Errorf("Expected empty result for empty content, got %d entries", len(result))
		}
	})

	// Test for handling for a non-existent word in TCI
	t.Run("WordsWithZeroCount", func(t *testing.T) {
		content := "apple banana apple"

		// Include a word in TCI that doesn't appear in content
		tci := TermCountIndex{
			"apple":  2,
			"banana": 1,
			"cherry": 0, // This word doesn't appear in content
		}

		result := getWordFrequency(content, wordSplitter, tci)

		// cherry should have frequency 0.0
		if freq, exists := result["cherry"]; !exists {
			t.Error("Expected 'cherry' to be present even with 0 count")
		} else if freq != 0.0 {
			t.Errorf("Expected frequency 0.0 for 'cherry', got %.4f", freq)
		}
	})

	// Test for handling punctuation in content
	t.Run("PunctuationInContent", func(t *testing.T) {
		content := "Hello, world! Hello again."

		tci := TermCountIndex{
			"hello": 2,
			"world": 1,
			"again": 1,
		}

		result := getWordFrequency(content, wordSplitter, tci)

		// Total words after splitting: ["hello", "", "world", "", "hello", "again", ""]
		// Length = 7, but some might be empty
		totalWords := len(wordSplitter.Split(content, -1))

		expectedHelloFreq := 2.0 / float64(totalWords)
		if freq, exists := result["hello"]; !exists {
			t.Error("Expected 'hello' to be present")
		} else if math.Abs(freq-expectedHelloFreq) > 0.0001 {
			t.Errorf("Expected frequency %.4f for 'hello', got %.4f", expectedHelloFreq, freq)
		}
	})

	// Test for handling case sensitivity
	t.Run("CaseHandling", func(t *testing.T) {
		content := "Apple APPLE apple"

		// TCI should have lowercase keys since getWordCount converts to lowercase
		tci := TermCountIndex{
			"apple": 3,
		}

		result := getWordFrequency(content, wordSplitter, tci)

		// All instances should be counted as "apple" (3 occurrences out of 3 total words)
		expectedFreq := 1.0
		if freq, exists := result["apple"]; !exists {
			t.Error("Expected 'apple' to be present")
		} else if math.Abs(freq-expectedFreq) > 0.0001 {
			t.Errorf("Expected frequency %.4f for 'apple', got %.4f", expectedFreq, freq)
		}

		// Uppercase versions should not exist in result
		if _, exists := result["Apple"]; exists {
			t.Error("Uppercase 'Apple' should not exist in result")
		}
		if _, exists := result["APPLE"]; exists {
			t.Error("Uppercase 'APPLE' should not exist in result")
		}
	})

	t.Run("LargeText", func(t *testing.T) {
		content := "the quick brown fox jumps over the lazy dog the fox is quick"

		tci := TermCountIndex{
			"the":   3,
			"quick": 2,
			"brown": 1,
			"fox":   2,
			"jumps": 1,
			"over":  1,
			"lazy":  1,
			"dog":   1,
			"is":    1,
		}

		result := getWordFrequency(content, wordSplitter, tci)

		// Check that all frequencies sum up correctly
		totalWords := len(wordSplitter.Split(content, -1))

		// Verify some specific frequencies
		expectedTheFreq := 3.0 / float64(totalWords)
		if freq, exists := result["the"]; !exists {
			t.Error("Expected 'the' to be present")
		} else if math.Abs(freq-expectedTheFreq) > 0.0001 {
			t.Errorf("Expected frequency %.4f for 'the', got %.4f", expectedTheFreq, freq)
		}

		expectedQuickFreq := 2.0 / float64(totalWords)
		if freq, exists := result["quick"]; !exists {
			t.Error("Expected 'quick' to be present")
		} else if math.Abs(freq-expectedQuickFreq) > 0.0001 {
			t.Errorf("Expected frequency %.4f for 'quick', got %.4f", expectedQuickFreq, freq)
		}
	})

	t.Run("CustomWordSplitter", func(t *testing.T) {
		content := "word1-word2_word3 word4"
		customSplitter := regexp.MustCompile(`[\s\-_]+`)

		tci := TermCountIndex{
			"word1": 1,
			"word2": 1,
			"word3": 1,
			"word4": 1,
		}

		result := getWordFrequency(content, customSplitter, tci)

		// With custom splitter, should split into 4 words
		totalWords := len(customSplitter.Split(content, -1))
		expectedFreq := 1.0 / float64(totalWords)

		words := []string{"word1", "word2", "word3", "word4"}
		for _, word := range words {
			if freq, exists := result[word]; !exists {
				t.Errorf("Expected '%s' to be present with custom splitter", word)
			} else if math.Abs(freq-expectedFreq) > 0.0001 {
				t.Errorf("Expected frequency %.4f for '%s', got %.4f", expectedFreq, word, freq)
			}
		}
	})

	t.Run("MismatchedTCIAndContent", func(t *testing.T) {
		// TCI has words that don't appear in content
		content := "apple banana"

		tci := TermCountIndex{
			"apple":  1,
			"banana": 1,
			"cherry": 1, // This word is not in content
			"date":   2, // Neither is this
		}

		result := getWordFrequency(content, wordSplitter, tci)

		// Words not in content should have frequency 0
		if freq, exists := result["cherry"]; !exists {
			t.Error("Expected 'cherry' to be present")
		} else if freq != 0.0 {
			t.Errorf("Expected frequency 0.0 for 'cherry' (not in content), got %.4f", freq)
		}

		if freq, exists := result["date"]; !exists {
			t.Error("Expected 'date' to be present")
		} else if freq != 0.0 {
			t.Errorf("Expected frequency 0.0 for 'date' (not in content), got %.4f", freq)
		}
	})
}

// Benchmark test for getWordFrequency function
func BenchmarkGetWordFrequency(b *testing.B) {
	content := "The quick brown fox jumps over the lazy dog. " +
		"This is a sample text for benchmarking the word frequency function. " +
		"It contains various words that should be processed efficiently."

	wordSplitter := regexp.MustCompile(`\W+`)

	// Create a realistic TCI
	tci := TermCountIndex{
		"the":          3,
		"quick":        1,
		"brown":        1,
		"fox":          1,
		"jumps":        1,
		"over":         1,
		"lazy":         1,
		"dog":          1,
		"this":         1,
		"is":           1,
		"sample":       1,
		"text":         1,
		"for":          1,
		"benchmarking": 1,
		"word":         1,
		"frequency":    1,
		"function":     1,
		"it":           1,
		"contains":     1,
		"various":      1,
		"words":        1,
		"that":         1,
		"should":       1,
		"be":           1,
		"processed":    1,
		"efficiently":  1,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		getWordFrequency(content, wordSplitter, tci)
	}
}
