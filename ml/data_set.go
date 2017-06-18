package ml

import (
	"os"
	"log"
	"bufio"
	"strings"
	"strconv"
	"math"
	"math/rand"
)

type DataSet [][]string

func FromCSV(filename string) DataSet {
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

	return DataSet(lines)
}

func (dataSet DataSet) StringValues(column int) []string {
	var values []string

	for i := range dataSet {
		values = append(values, dataSet[i][column])
	}

	return values
}

func (dataSet DataSet) Float64Values(column int) []float64 {
	var values []float64

	for i := range dataSet {
		var value, err = strconv.ParseFloat(dataSet[i][column], 64)

		if err != nil {
			log.Fatalf("Cannot parse float number %s\n", dataSet[i][column])
		}

		values = append(values, value)
	}

	return values
}

func (dataSet DataSet) IndexedValues(column int) ([]int, map[int]string) {
	var namesToIndexes = make(map[string]int)
	var indexesToNames = make(map[int]string)

	var index = 0
	for i := range dataSet {
		var name = dataSet[i][column]

		if _, exists := namesToIndexes[name]; !exists {
			namesToIndexes[name] = index
			indexesToNames[index] = name
			index += 1
		}
	}

	var indexes []int

	for i := range dataSet {
		var index = namesToIndexes[dataSet[i][column]]
		indexes = append(indexes, index)
	}

	return indexes, indexesToNames
}

func (dataSet DataSet) MinMax(column int) (float64, float64) {
	var values = dataSet.Float64Values(column)

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

func (dataSet DataSet) NormalizedFloat64Values(column int) []float64 {
	var values = dataSet.Float64Values(column)

	var normalizedFloat64Values []float64

	var min, max = dataSet.MinMax(column)

	for row := range dataSet {
		var value = values[row]
		normalizedFloat64Values = append(normalizedFloat64Values, (value - min) / (max - min))
	}

	return normalizedFloat64Values
}

func (dataSet DataSet) Mean(column int) float64 {
	var sum = float64(0)

	var values = dataSet.Float64Values(column)

	for _, value := range values {
		sum += value
	}

	return sum / float64(len(dataSet))
}

func (dataSet DataSet) StandardDeviation(column int) float64 {
	var mean = dataSet.Mean(column)

	var variance []float64

	for _, value := range dataSet.Float64Values(column) {
		variance = append(variance, math.Pow(value - mean, 2))
	}

	var standardDeviation = float64(0)

	for _, v := range variance {
		standardDeviation += v
	}

	standardDeviation = math.Sqrt(standardDeviation / float64(len(dataSet) - 1))

	return standardDeviation
}

func (dataSet DataSet) StandardizedFloat64Values(column int) []float64 {
	var values = dataSet.Float64Values(column)

	var standardizedFloat64Values []float64

	var mean = dataSet.Mean(column)
	var standardDeviation = dataSet.StandardDeviation(column)

	for row := range dataSet {
		var value = values[row]
		standardizedFloat64Values = append(standardizedFloat64Values, (value - mean) / standardDeviation)
	}

	return standardizedFloat64Values
}

func (dataSet DataSet) TrainTestSplit(split float64) (DataSet, DataSet) {
	var train DataSet

	var trainSize = int(split * float64(len(dataSet)))
	var test = dataSet

	for len(train) < trainSize {
		var index = rand.Intn(len(test))

		var row = test[index]

		test = append(test[:index], test[index + 1:]...)

		train = append(train, row)
	}

	return train, test
}

func (dataSet DataSet) CrossValidationSplit(folds int) []DataSet {
	var dataSetSplit []DataSet

	var dataSetCopy = dataSet

	var foldSize = int(len(dataSet) / folds)

	for i := 0; i < folds; i++ {
		var fold DataSet

		for len(fold) < foldSize {
			var index = rand.Intn(len(dataSetCopy))

			var row = dataSetCopy[index]

			dataSetCopy = append(dataSetCopy[:index], dataSetCopy[index + 1: ]...)

			fold = append(fold, row)
		}

		dataSetSplit = append(dataSetSplit, DataSet(fold))
	}

	return dataSetSplit
}

func AccuracyMetric(actual []string, predicted []string) float64 {
	var correct = float64(0)

	for i := 0; i < len(actual); i++ {
		if actual[i] == predicted[i] {
			correct += 1
		}
	}

	return correct / float64(len(actual)) * 100.0
}

func ConfusionMatrix(actual []string, predicted []string) [][]int64 {
	var unique = actual

	var matrix [][]int64

	for i := 0; i < len(unique); i++ {
		matrix = append(matrix, []int64{})
	}

	for i := 0; i < len(unique); i++ {
		for j := 0; j< len(unique); j++ {
			matrix[i] = append(matrix[i], 0)
		}
	}

	return matrix
}

