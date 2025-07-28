package main

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

func LoadStopwords(filePath string) (map[string]struct{}, error) {
	stopwords := make(map[string]struct{})
	file, err := os.Open(filePath)
	if err != nil {
		return nil, errors.New("Error opening file, check file path: " + err.Error())
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		word := strings.TrimSpace(scanner.Text())
		stopwords[word] = struct{}{} // add the word to the stopwords map
		delete(stopwords, "")        // remove empty strings if any
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if len(stopwords) == 0 {
		return nil, errors.New("no stopwords found in file")
	}
	return stopwords, nil
}
