package simutil

import (
	"encoding/json"
	"fmt"
	"os"
)

func WriteJsonFile(obj interface{}, outputDirectory string, outputJsonFileName string) {
	if err := os.MkdirAll(outputDirectory, os.ModePerm); err != nil {
		panic(fmt.Sprintf("Cannot create output directory (%s)", err))
	}

	file, err := os.Create(outputDirectory + "/" + outputJsonFileName)

	if err != nil {
		panic(fmt.Sprintf("Cannot create JSON file (%s)", err))
	}

	defer file.Close()

	j, err := json.MarshalIndent(obj, "", "  ")

	if err != nil {
		panic(fmt.Sprintf("Cannot encode object to JSON (%s)", err))
	}

	if _, err := file.Write(j); err != nil {
		panic(fmt.Sprintf("Cannot write JSON file (%s)", err))
	}
}

func LoadJsonFile(outputDirectory string, outputJsonFileName string, data interface{}) {
	var file, err = os.Open(outputDirectory + "/" + outputJsonFileName)

	if err != nil {
		panic(fmt.Sprintf("Cannot open JSON file (%s)", err))
	}

	defer file.Close()

	var jsonParser = json.NewDecoder(file)

	if err := jsonParser.Decode(data); err != nil {
		panic(fmt.Sprintf("Cannot decode object from JSON (%s)", err))
	}
}
