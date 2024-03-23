package main

import (
	"bufio"
	"encoding/csv"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

func loadRandomNumberOfLines(rows []map[string]string) []map[string]string {
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)

	numRandomLines := ConfigInstance.NumberOfSamples
	randomIndices := r.Perm(len(rows))

	// Select the random lines
	var randomLines []map[string]string
	for i := 0; i < numRandomLines && i < len(rows); i++ {
		randomLines = append(randomLines, rows[randomIndices[i]])
	}

	return randomLines
}

func readMotivationFile() []map[string]string {
	filePath := ConfigInstance.ProducerDataFile

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("os.Open: %s", err)
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatalf("file.Close: %s", err)
		}
	}(file)

	reader := csv.NewReader(bufio.NewReader(file))
	header, err := reader.Read()
	if err != nil {
		log.Fatalf("reader.Read: %s", err)
	}

	for i, h := range header {
		header[i] = strings.ToLower(h)
	}

	var rows []map[string]string
	for {
		row, err := reader.Read()
		if err != nil {
			break
		}

		m := make(map[string]string)
		for i, value := range row {
			m[header[i]] = value
		}
		rows = append(rows, m)
	}

	return loadRandomNumberOfLines(rows)
}
