package main

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

// takes a filepath to a list of stopwords in a text file and returns a map of stopwords
func LoadStopwords(filePath string) (map[string]struct{}, error) {
	// make empty map
	stopwords := make(map[string]struct{})

	// open the file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, errors.New("Error opening file, check file path: " + err.Error())
	}
	defer file.Close()

	// create a new scanner and for each line in the file, add the word to the stopwords map
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		word := strings.TrimSpace(scanner.Text())
		stopwords[word] = struct{}{} // add the word to the stopwords map
		delete(stopwords, "")        // remove empty strings if any
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// if there are no stopwords in the file, return an error
	if len(stopwords) == 0 {
		return nil, errors.New("no stopwords found in file")
	}
	return stopwords, nil
}
