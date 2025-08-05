package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

/*
This file tests for:
- correct keyword extraction
- not enough keywords to return
- empty TermFrequencyIndex handling
- zero topN handling
- single keyword handling
- equal frequencies handling
- order of keywords
- negative topN handling

*/

func TestGetKeywords(t *testing.T) {
	// Helper function to capture stdout
	captureOutput := func(f func()) string {
		old := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		f()

		w.Close()
		os.Stdout = old

		var buf bytes.Buffer
		io.Copy(&buf, r)
		return buf.String()
	}

	t.Run("BasicKeywordExtraction", func(t *testing.T) {
		wordFreq := TermFrequencyIndex{
			"programming": 0.3,
			"computer":    0.2,
			"science":     0.15,
			"algorithm":   0.1,
			"data":        0.05,
		}

		output := captureOutput(func() {
			getKeywords(wordFreq, 3)
		})

		lines := strings.Split(strings.TrimSpace(output), "\n")

		// Should have exactly 3 lines (top 3 keywords)
		if len(lines) != 3 {
			t.Errorf("Expected 3 keywords, got %d", len(lines))
		}

		// Check if keywords are in descending order of frequency
		expectedOrder := []string{"programming", "computer", "science"}
		for i, expected := range expectedOrder {
			if i < len(lines) && !strings.Contains(lines[i], expected) {
				t.Errorf("Expected keyword '%s' at position %d, but got '%s'", expected, i+1, lines[i])
			}
		}
	})

	// test for when there are not enough keywords to return
	t.Run("RequestMoreThanAvailable", func(t *testing.T) {
		wordFreq := TermFrequencyIndex{
			"apple":  0.5,
			"banana": 0.3,
		}

		output := captureOutput(func() {
			getKeywords(wordFreq, 5) // Request 5 but only 2 available
		})

		lines := strings.Split(strings.TrimSpace(output), "\n")

		// Should only return 2 keywords (all available)
		if len(lines) != 2 {
			t.Errorf("Expected 2 keywords (all available), got %d", len(lines))
		}
	})

	// handle an empty TermFrequencyIndex
	t.Run("EmptyWordFrequency", func(t *testing.T) {
		wordFreq := TermFrequencyIndex{}

		output := captureOutput(func() {
			getKeywords(wordFreq, 3)
		})

		// Should produce no output for empty frequency map
		if strings.TrimSpace(output) != "" {
			t.Errorf("Expected no output for empty frequency map, got: '%s'", output)
		}
	})

	// Test for zero topN
	t.Run("ZeroTopN", func(t *testing.T) {
		wordFreq := TermFrequencyIndex{
			"word1": 0.5,
			"word2": 0.3,
		}

		output := captureOutput(func() {
			getKeywords(wordFreq, 0)
		})

		// Should produce no output when topN is 0
		if strings.TrimSpace(output) != "" {
			t.Errorf("Expected no output when topN is 0, got: '%s'", output)
		}
	})

	// Test for a single keyword
	t.Run("SingleKeyword", func(t *testing.T) {
		wordFreq := TermFrequencyIndex{
			"singleton": 1.0,
		}

		output := captureOutput(func() {
			getKeywords(wordFreq, 1)
		})

		lines := strings.Split(strings.TrimSpace(output), "\n")

		if len(lines) != 1 {
			t.Errorf("Expected 1 keyword, got %d", len(lines))
		}

		if !strings.Contains(lines[0], "singleton") {
			t.Errorf("Expected 'singleton' in output, got: '%s'", lines[0])
		}
	})

	// Test for equal frequencies
	t.Run("EqualFrequencies", func(t *testing.T) {
		wordFreq := TermFrequencyIndex{
			"word1": 0.25,
			"word2": 0.25,
			"word3": 0.25,
			"word4": 0.25,
		}

		output := captureOutput(func() {
			getKeywords(wordFreq, 2)
		})

		lines := strings.Split(strings.TrimSpace(output), "\n")

		// Should return exactly 2 keywords even with equal frequencies
		if len(lines) != 2 {
			t.Errorf("Expected 2 keywords, got %d", len(lines))
		}

		// All returned keywords should be from our original set
		validWords := map[string]bool{
			"word1": true,
			"word2": true,
			"word3": true,
			"word4": true,
		}

		for _, line := range lines {
			found := false
			for word := range validWords {
				if strings.Contains(line, word) {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Unexpected keyword in output: '%s'", line)
			}
		}
	})

	// Test order of keywords
	t.Run("MixedFrequencies", func(t *testing.T) {
		wordFreq := TermFrequencyIndex{
			"high":   0.4,
			"medium": 0.3,
			"low":    0.2,
			"lower":  0.1,
		}

		output := captureOutput(func() {
			getKeywords(wordFreq, 4)
		})

		lines := strings.Split(strings.TrimSpace(output), "\n")

		if len(lines) != 4 {
			t.Errorf("Expected 4 keywords, got %d", len(lines))
		}

		// Check the order is correct (highest frequency first)
		expectedOrder := []string{"high", "medium", "low", "lower"}
		for i, expected := range expectedOrder {
			if i < len(lines) && !strings.Contains(lines[i], expected) {
				t.Errorf("Expected keyword '%s' at position %d, but line was '%s'", expected, i+1, lines[i])
			}
		}
	})

	// Test for negative topN
	t.Run("NegativeTopN", func(t *testing.T) {
		wordFreq := TermFrequencyIndex{
			"word": 0.5,
		}

		output := captureOutput(func() {
			getKeywords(wordFreq, -1)
		})

		// Should produce no output for negative topN
		if strings.TrimSpace(output) != "" {
			t.Errorf("Expected no output for negative topN, got: '%s'", output)
		}
	})
}

// Benchmark test for getKeywords function
func BenchmarkGetKeywords(b *testing.B) {
	// Create a sample word frequency map
	wordFreq := TermFrequencyIndex{
		"programming": 0.3,
		"computer":    0.25,
		"science":     0.2,
		"algorithm":   0.15,
		"data":        0.1,
		"structure":   0.08,
		"analysis":    0.06,
		"software":    0.04,
		"development": 0.02,
		"technology":  0.01,
	}

	// Redirect output to discard it during benchmarking
	old := os.Stdout
	os.Stdout, _ = os.Open(os.DevNull)
	defer func() { os.Stdout = old }()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		getKeywords(wordFreq, 5)
	}
}
