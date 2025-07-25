package main

import "sort"

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
		println(sorted[i].Key, ":", sorted[i].Value)
	}
}
