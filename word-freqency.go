package main

import (
	"errors"
	"regexp"
	"sort"
	"strings"
	"unicode"
)

type TermCountIndex map[string]int
type TermFrequencyIndex map[string]float64

// Returns a map of words and their counts from the content, excluding stopwords and short words
// Takes a content string, a map of stopwords, and a regex for splitting words
func getWordCount(content string, stopwords map[string]struct{}, wordSplitter *regexp.Regexp) (TermCountIndex, error) {
	tci := make(TermCountIndex)

	// get words
	words := wordSplitter.Split(strings.ToLower(content), -1)
	for _, word := range words {
		// if word is not a space or digit and is longer than two and not a stopword
		if word != "" && !unicode.IsDigit(rune(word[0])) && len(word) > 2 {
			if _, ok := stopwords[word]; !ok {
				// add word to tfi and increase count
				tci[word]++
			}
		}
	}
	// handle no valid words case
	if len(tci) == 0 {
		return nil, errors.New("No valid words found in content")
	}

	return tci, nil
}

// Calculates the term frequency index from the content and word count by dividing the number of times each word appears by the total number of words
// Takes a content string, a regex for splitting words, and a TermCountIndex which is made from getWordCount
func getWordFrequency(content string, wordSplitter *regexp.Regexp, tci TermCountIndex) TermFrequencyIndex {
	// make tfi
	tfi := make(TermFrequencyIndex)

	// get words
	words := wordSplitter.Split(strings.ToLower(content), -1)
	totalWords := len(words)

	// for each word, divide the number of times it appears by the total number of words
	for _, word := range words {
		tfi[word] = float64(tci[word]) / float64(totalWords)
	}

	return tfi

}

// getKeywords takes a TermFrequencyIndex and returns the top N keywords based on which words are most frequent
func getKeywords(wordFrequency TermFrequencyIndex, topN int) {
	// Sort the word frequency index by frequency
	type kv struct {
		Key   string
		Value float64
	}
	var sorted []kv

	// Convert the map to a slice of key value pairs
	for k, v := range wordFrequency {
		sorted = append(sorted, kv{k, v})
	}

	// Sort by value (frequency)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Value > sorted[j].Value
	})

	// Print top N keywords
	for i := 0; i < topN && i < len(sorted); i++ {
		println(sorted[i].Key)
	}
}
