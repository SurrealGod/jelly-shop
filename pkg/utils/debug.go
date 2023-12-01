package utils

import (
	"fmt"

	"encoding/json"
)

func Debug(data any) {
	bytes, _ := json.MarshalIndent(data, "","\t")
	fmt.Println(string(bytes))
}

func Output(data any) []byte {
	bytes, _ := json.Marshal(data)
	return bytes
}