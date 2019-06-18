package ml

import (
	"fmt"
	"testing"
)

func TestFromCSV(t *testing.T) {
	var filename = "Data/iris.csv"
	var dataSet = FromCSV(filename)

	fmt.Printf("Loaded data file %s with %d rows and %d columns.\n", filename, len(dataSet), len(dataSet[0]))
	fmt.Println(dataSet[0])
	fmt.Println(dataSet.Float64Values(3))
	fmt.Println(dataSet.IndexedValues(len(dataSet[0]) - 1))
}

func TestDataSet_MinMax(t *testing.T) {
	var dataSet = DataSet{{"50", "30"}, {"20", "90"}}
	fmt.Println(dataSet)

	var min, max = dataSet.MinMax(0)
	fmt.Println(min, max)
}

func TestDataSet_NormalizedFloat64Values(t *testing.T) {
	var filename = "Data/pima-indians-diabetes.csv"
	var dataSet = FromCSV(filename)

	fmt.Println(dataSet.Float64Values(0))

	var min, max = dataSet.MinMax(0)
	fmt.Println(min, max)

	fmt.Println(dataSet.NormalizedFloat64Values(0))
}

func TestDataSet_Mean(t *testing.T) {
	var dataSet = DataSet{{"50", "30"}, {"20", "90"}, {"30", "50"}}
	fmt.Println(dataSet)

	var mean = dataSet.Mean(1)
	fmt.Println(mean)
}

func TestDataSet_StandardDeviation(t *testing.T) {
	var dataSet = DataSet{{"50", "30"}, {"20", "90"}, {"30", "50"}}
	fmt.Println(dataSet)

	var standardDeviation = dataSet.StandardDeviation(1)
	fmt.Println(standardDeviation)
}

func TestDataSet_StandardizedFloat64Values(t *testing.T) {
	var dataSet = DataSet{{"50", "30"}, {"20", "90"}, {"30", "50"}}
	fmt.Println(dataSet)

	fmt.Println(dataSet.Float64Values(0))

	var min, max = dataSet.MinMax(0)
	fmt.Println(min, max)

	fmt.Println(dataSet.StandardizedFloat64Values(0))
}

func TestDataSet_TrainTestSplit(t *testing.T) {
	var dataSet = DataSet{{"1"}, {"2"}, {"3"}, {"4"}, {"5"}, {"6"}, {"7"}, {"8"}, {"9"}, {"10"}}
	fmt.Println(dataSet)

	var train, test = dataSet.TrainTestSplit(0.6)

	fmt.Println(train)
	fmt.Println(test)
}

func TestDataSet_CrossValidationSplit(t *testing.T) {
	var dataSet = DataSet{{"1"}, {"2"}, {"3"}, {"4"}, {"5"}, {"6"}, {"7"}, {"8"}, {"9"}, {"10"}}
	fmt.Println(dataSet)

	var folds = dataSet.CrossValidationSplit(4)

	fmt.Println(folds)
}

func TestAccuracyMetric(t *testing.T) {
	var actual = []string{"0", "0", "0", "0", "0", "1", "1", "1", "1", "1"}
	var predicted = []string{"0", "1", "0", "0", "0", "1", "0", "1", "1", "1"}

	var accuracy = AccuracyMetric(actual, predicted)

	fmt.Println(accuracy)
}
