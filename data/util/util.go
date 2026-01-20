package util

import (
	"encoding/csv"
	"fmt"
	"os"
)

func ReadFile(fileName string) ([][]string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("could not open file %s: %w", fileName, err)
	}
	defer file.Close()

	lines, err := csv.NewReader(file).ReadAll()
	if err != nil {
		return nil, fmt.Errorf("could not read file %s: %w", fileName, err)
	}

	return lines[1:], nil
}
