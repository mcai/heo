package simutil

import (
	"encoding/json"
	"bytes"
)

const STATS_JSON_FILE_NAME = "stats.json"

type Stat struct {
	Key   string
	Value interface{}
}

type Stats []Stat

func (stats Stats) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer

	buf.WriteString("{")

	for i, stat := range stats {
		if i != 0 {
			buf.WriteString(",")
		}

		key, err := json.Marshal(stat.Key)
		if err != nil {
			return nil, err
		}
		buf.Write(key)

		buf.WriteString(":")

		val, err := json.Marshal(stat.Value)
		if err != nil {
			return nil, err
		}
		buf.Write(val)
	}

	buf.WriteString("}")

	return buf.Bytes(), nil
}
