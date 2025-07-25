package main

import (
	"regexp"
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
