package main

import (
	"regexp"
	"sort"
	"strings"
	"unicode"
)

type TermCountIndex map[string]int
type TermFrequencyIndex map[string]float64

func getWordCount(content string, stopwords map[string]struct{}, wordSplitter *regexp.Regexp) TermCountIndex {
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

	// fmt.Printf("Term count Index: %v\n", tci)
	return tci
}

func getWordFrequency(content string, stopwords map[string]struct{}, wordSplitter *regexp.Regexp, tci TermCountIndex) TermFrequencyIndex {
	// make tfi
	tfi := make(TermFrequencyIndex)

	// get words
	words := wordSplitter.Split(strings.ToLower(content), -1)
	totalWords := len(words)

	// for each word
	for _, word := range words {
		tfi[word] = float64(tci[word]) / float64(totalWords)
	}

	// fmt.Printf("Term Frequency Index: %v\n", tfi)
	return tfi

}

func getKeywords(wordFrequency TermFrequencyIndex, topN int) {
	// Sort the word frequency index by frequency
	type kv struct {
		Key   string
		Value float64
	}
	var sorted []kv
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
