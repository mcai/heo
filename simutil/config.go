package simutil

const (
	CPU_CONFIG_JSON_FILE_NAME = "config.cpu.json"
	UNCORE_CONFIG_JSON_FILE_NAME = "config.uncore.json"
	NOC_CONFIG_JSON_FILE_NAME = "config.noc.json"
)

type Config interface {
	Dump(outputDirectory string)
}
