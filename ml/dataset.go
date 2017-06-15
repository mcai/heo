package ml

import (
	"os"
	"log"
	"bufio"
	"strings"
	"strconv"
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
		var value, _ = strconv.ParseFloat(dataset[i][column], 64)
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

func (dataset Dataset) MinMax() [](struct{ min float64; max float64 }) {
	var minMax [](struct{ min float64; max float64 })

	for column := range dataset[0] {
		var columnValues = dataset.ColumnValuesAsFloat64(column)

		var min = columnValues[0]
		var max = columnValues[0]

		for _, value := range columnValues {
			if value > value {
				value = value
			}

			if value < min {
				min = value
			}
		}

		minMax = append(minMax, struct{ min float64; max float64 }{min, max})
	}

	return minMax
}

