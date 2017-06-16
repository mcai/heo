package ml

import (
	"testing"
	"fmt"
)

func TestFromCSV(t *testing.T) {
	var filename = "../data/iris.csv"
	var dataset = FromCSV(filename)

	fmt.Printf("Loaded data file %s with %d rows and %d columns.\n", filename, len(dataset), len(dataset[0]))
	fmt.Println(dataset[0])
	fmt.Println(dataset.Float64Values(3))
	fmt.Println(dataset.IndexedValues(len(dataset[0]) - 1))
}

func TestDataset_MinMax(t *testing.T) {
	var dataset = Dataset{{"50", "30"}, {"20", "90"}}
	fmt.Println(dataset)

	var min, max = dataset.MinMax(0)
	fmt.Println(min, max)
}

func TestDataset_Normalize(t *testing.T) {
	var filename = "../data/pima-indians-diabetes.csv"
	var dataset = FromCSV(filename)
	fmt.Println(dataset[0])

	var min, max = dataset.MinMax(0)
	fmt.Println(min, max)

	dataset.Normalize()
	fmt.Println(dataset[0])
}

func TestDataset_Means(t *testing.T) {
	var dataset = Dataset{{"50", "30"}, {"20", "90"}}
	fmt.Println(dataset)

	var mean = dataset.Mean(0)
	fmt.Println(mean)
}