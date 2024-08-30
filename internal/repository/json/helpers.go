package json

import (
	"encoding/json"
	"io"
	"os"
)

func initFile(fileName string, v any) error {

	bytes, err := readData(fileName)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, v)
}

func readData(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}
