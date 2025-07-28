package main

import (
	"os"
	"regexp"
)

func getStringKeywords(content string, numKeywords int) {
	// parse text
	stopwords, err := LoadStopwords("./data/stopwords.txt")
	if err != nil {
		print("Error loading stopwords:", err)
	}

	wordSplitter := regexp.MustCompile(`[^a-zA-Z0-9]+`)

	wordCount := getWordCount(content, stopwords, wordSplitter)
	wordFrequency := getWordFrequency(content, stopwords, wordSplitter, wordCount)

	getKeywords(wordFrequency, numKeywords)

}

func getFileKeywords(filePath string, numKeywords int) {
	// parse file
	content, err := LoadFileContent(filePath)
	if err != nil {
		print("Error loading file content:", err)
		return
	}

	stopwords, err := LoadStopwords("./data/stopwords.txt")
	if err != nil {
		print("Error loading stopwords:", err)
		return
	}

	wordSplitter := regexp.MustCompile(`[^a-zA-Z0-9]+`)

	wordCount := getWordCount(content, stopwords, wordSplitter)
	wordFrequency := getWordFrequency(content, stopwords, wordSplitter, wordCount)

	getKeywords(wordFrequency, numKeywords)
}

func LoadFileContent(file string) (string, error) {
	content, err := os.ReadFile(file)
	if err != nil {
		return "", err
	}
	return string(content), nil
}
