package ml

import (
	"testing"
	"fmt"
)

func TestLoadCSV(t *testing.T) {
	var filename = "../data/iris.csv"
	var dataset = FromCSV(filename)

	fmt.Printf("Loaded data file %s with %d rows and %d columns.\n", filename, len(dataset), len(dataset[0]))
	fmt.Println(dataset[0])
	fmt.Println(dataset.ColumnValuesAsFloat64(4))
	fmt.Println(dataset.ColumnValuesAsIntMap(len(dataset[0]) - 1))
}

func TestDatasetMinMax(t *testing.T) {
	var dataset = Dataset{{"50", "30"}, {"20", "90"}}
	fmt.Println(dataset)
	var minMax = dataset.MinMax()
	fmt.Println(minMax)
}