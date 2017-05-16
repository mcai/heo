package simutil

import (
	"fmt"
	"time"
)

type Experiment interface {
	Run(skipIfStatsFileExists bool)
}

func RunExperiments(experiments []Experiment, skipIfStatsFileExists bool) {
	var done = make(chan bool)

	for i, e := range experiments {
		go func(i int, experiment Experiment, c chan bool) {
			var l = len(experiments)

			fmt.Printf("[%s] Experiment %d/%d started.\n",
				time.Now().Format("2006-01-02 15:04:05"), i + 1, l)

			experiment.Run(skipIfStatsFileExists)

			done <- true

			fmt.Printf("[%s] Experiment %d/%d ended.\n",
				time.Now().Format("2006-01-02 15:04:05"), i + 1, l)
		}(i, e, done)
	}

	for i := 0; i < len(experiments); i++ {
		<-done

		fmt.Printf("[%s] There are %d experiments to be run.\n",
			time.Now().Format("2006-01-02 15:04:05"), len(experiments) - i - 1)
	}
}
