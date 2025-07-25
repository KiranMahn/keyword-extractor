package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func LoadStopwords(filePath string) (map[string]struct{}, error) {
	stopwords := make(map[string]struct{})
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("err:", err) // Debugging print
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		word := strings.TrimSpace(scanner.Text())
		stopwords[word] = struct{}{} // add the word to the stopwords map
		// fmt.Println("Loaded stopword:", word) // Debugging print
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return stopwords, nil
}
