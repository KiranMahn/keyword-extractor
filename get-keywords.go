package main

import (
	"os"
	"regexp"
)

// finds keywwords for text in a string
func getStringKeywords(content string, numKeywords int) []string {
	// parse text
	stopwords, err := LoadStopwords("./data/stopwords.txt")
	if err != nil {
		print("Error loading stopwords:", err)
	}
	wordSplitter := regexp.MustCompile(`[^a-zA-Z0-9]+`)

	// get word count and frequency
	wordCount, err := getWordCount(content, stopwords, wordSplitter)
	if err != nil {
		print("Error getting word count:", err)
		return nil
	}
	wordFrequency := getWordFrequency(content, wordSplitter, wordCount)

	// get keywords based on frquency of words that are not stopwords
	keywords := getKeywords(wordFrequency, numKeywords)

	return keywords

}

// finds keywords for text from a given filepath
func getFileKeywords(filePath string, numKeywords int) []string {
	// load file to string
	content, err := LoadFileContent(filePath)
	if err != nil {
		print("Error loading file content:", err)
		return nil
	}

	// get keywords from string content
	keywords := getStringKeywords(content, numKeywords)

	return keywords
}

// LoadFileContent reads the content of a file and returns it as a string.
func LoadFileContent(file string) (string, error) {
	content, err := os.ReadFile(file)
	if err != nil {
		return "", err
	}
	return string(content), nil
}
