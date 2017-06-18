package ml

import (
	"testing"
	"fmt"
	//"math/rand"
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

func TestDataset_NormalizedFloat64Values(t *testing.T) {
	var filename = "../data/pima-indians-diabetes.csv"
	var dataset = FromCSV(filename)

	fmt.Println(dataset.Float64Values(0))

	var min, max = dataset.MinMax(0)
	fmt.Println(min, max)

	fmt.Println(dataset.NormalizedFloat64Values(0))
}

func TestDataset_Means(t *testing.T) {
	var dataset = Dataset{{"50", "30"}, {"20", "90"}, {"30", "50"}}
	fmt.Println(dataset)

	var mean = dataset.Mean(1)
	fmt.Println(mean)
}

func TestDataset_Stdev(t *testing.T) {
	var dataset = Dataset{{"50", "30"}, {"20", "90"}, {"30", "50"}}
	fmt.Println(dataset)

	var stdev = dataset.Stdev(1)
	fmt.Println(stdev)
}

func TestDataset_StandardizedFloat64Values(t *testing.T) {
	var dataset = Dataset{{"50", "30"}, {"20", "90"}, {"30", "50"}}
	fmt.Println(dataset)

	fmt.Println(dataset.Float64Values(0))

	var min, max = dataset.MinMax(0)
	fmt.Println(min, max)

	fmt.Println(dataset.StandardizedFloat64Values(0))
}

func TestDataset_TrainTestSplit(t *testing.T) {
	var dataset = Dataset{{"1"}, {"2"}, {"3"}, {"4"}, {"5"}, {"6"}, {"7"}, {"8"}, {"9"}, {"10"}}
	fmt.Println(dataset)

	var train, test = dataset.TrainTestSplit(0.6)

	fmt.Println(train)
	fmt.Println(test)
}

func TestDataset_CrossValidationSplit(t *testing.T) {
	var dataset = Dataset{{"1"}, {"2"}, {"3"}, {"4"}, {"5"}, {"6"}, {"7"}, {"8"}, {"9"}, {"10"}}
	fmt.Println(dataset)

	var folds = dataset.CrossValidationSplit(4)

	fmt.Println(folds)
}

func TestAccuracyMetric(t *testing.T) {
	var actual = []string{"0", "0", "0", "0", "0", "1", "1", "1", "1", "1"}
	var predicted = []string{"0", "1", "0", "0", "0", "1", "0", "1", "1", "1"}

	var accuracy = AccuracyMetric(actual, predicted)

	fmt.Println(accuracy)
}