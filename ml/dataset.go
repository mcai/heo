package ml

import (
	"os"
	"log"
	"bufio"
	"strings"
	"strconv"
)

func LoadCSV(filename string) [][]string {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal("Cannot open file: " + filename)
	}
	defer file.Close()

	var lines [][]string
	var scanner = bufio.NewScanner(file)
	for scanner.Scan() {
		var text = scanner.Text()
		var line = strings.Split(text, ",")

		if len(strings.TrimSpace(text)) > 0 {
			lines = append(lines, line)
		}
	}

	return lines
}

func StringColumnValuesToFloat64s(dataset [][]string, column int) []float64 {
	var values []float64

	for i, _ := range dataset {
		var value, _ = strconv.ParseFloat(dataset[i][column], 64)
		values = append(values, value)
	}

	return values
}

func StringColumnValuesToInts(dataset [][]string, column int) ([]int, map[int]string) {
	var class_values = make(map[string]int)
	var lookup = make(map[int]string)

	var index = 0
	for i, _ := range dataset {
		var value = dataset[i][column]

		if _, exists := class_values[value]; !exists {
			class_values[value] = index
			lookup[index] = value
			index += 1
		}
	}

	var values []int

	for i, _ := range dataset {
		var value = class_values[dataset[i][column]]
		values = append(values, value)
	}

	return values, lookup
}

