package simutil

type Experiment interface {
	Run(skipIfStatsFileExists bool)
}
