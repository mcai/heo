package ml

import (
	"os"
	"log"
	"bufio"
	"strings"
	"strconv"
	"fmt"
)

type Dataset [][]string

func FromCSV(filename string) Dataset {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal("Cannot open file: " + filename)
	}
	defer file.Close()

	var lines [][]string
	var scanner = bufio.NewScanner(file)
	for scanner.Scan() {
		var text = scanner.Text()

		if len(strings.TrimSpace(text)) > 0 {
			var line = strings.Split(text, ",")
			lines = append(lines, line)
		}
	}

	return Dataset(lines)
}

func (dataset Dataset) StringValues(column int) []string {
	var values []string

	for i := range dataset {
		values = append(values, dataset[i][column])
	}

	return values
}

func (dataset Dataset) Float64Values(column int) []float64 {
	var values []float64

	for i := range dataset {
		var value, err = strconv.ParseFloat(dataset[i][column], 64)

		if err != nil {
			log.Fatalf("Cannot parse float number %s\n", dataset[i][column])
		}

		values = append(values, value)
	}

	return values
}

func (dataset Dataset) IndexedValues(column int) ([]int, map[int]string) {
	var namesToIndexes = make(map[string]int)
	var indexesToNames = make(map[int]string)

	var index = 0
	for i := range dataset {
		var name = dataset[i][column]

		if _, exists := namesToIndexes[name]; !exists {
			namesToIndexes[name] = index
			indexesToNames[index] = name
			index += 1
		}
	}

	var indexes []int

	for i := range dataset {
		var index = namesToIndexes[dataset[i][column]]
		indexes = append(indexes, index)
	}

	return indexes, indexesToNames
}

func (dataset Dataset) MinMax(column int) (float64, float64) {
	var values = dataset.Float64Values(column)

	var min = values[0]
	var max = values[0]

	for _, value := range values {
		if value > max {
			max = value
		}

		if value < min {
			min = value
		}
	}

	return min, max
}

func (dataset Dataset) Normalize() {
	for _, row := range dataset {
		for i := 0; i < len(row); i++ {
			var value, err = strconv.ParseFloat(row[i], 64)

			if err == nil {
				var min, max = dataset.MinMax(i)
				row[i] = fmt.Sprintf("%f", (value - min) / (max - min))
			}
		}
	}
}

func (dataset Dataset) Mean(column int) float64 {
	var sum = float64(0)

	var values = dataset.Float64Values(column)

	for _, value := range values {
		sum += value
	}

	return sum / float64(len(dataset))
}

