package ml

import (
	"testing"
	"fmt"
)

func TestLoadCSV(t *testing.T) {
	var filename = "../data/iris.csv"
	var dataset = LoadCSV(filename)

	fmt.Printf("Loaded data file %s with %d rows and %d columns.\n", filename, len(dataset), len(dataset[0]))
	fmt.Println(dataset[0])
	fmt.Println(StringColumnValuesToFloat64s(dataset, 4))
	fmt.Println(StringColumnValuesToInts(dataset, len(dataset[0]) -1))
}
