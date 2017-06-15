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
		var line = strings.Split(text, ",")

		if len(strings.TrimSpace(text)) > 0 {
			lines = append(lines, line)
		}
	}

	return Dataset(lines)
}

func (dataset Dataset) ColumnValuesAsString(column int) []string {
	var values []string

	for i := range dataset {
		values = append(values, dataset[i][column])
	}

	return values
}

func (dataset Dataset) ColumnValuesAsFloat64(column int) []float64 {
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

func (dataset Dataset) ColumnValuesAsIntMap(column int) ([]int, map[int]string) {
	var classValues = make(map[string]int)
	var lookup = make(map[int]string)

	var index = 0
	for i := range dataset {
		var value = dataset[i][column]

		if _, exists := classValues[value]; !exists {
			classValues[value] = index
			lookup[index] = value
			index += 1
		}
	}

	var values []int

	for i := range dataset {
		var value = classValues[dataset[i][column]]
		values = append(values, value)
	}

	return values, lookup
}

func (dataset Dataset) MinMax(column int) (float64, float64) {
	var columnValues = dataset.ColumnValuesAsFloat64(column)

	var min = columnValues[0]
	var max = columnValues[0]

	for _, value := range columnValues {
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

