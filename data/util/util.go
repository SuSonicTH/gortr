package util

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
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

func Normalize(number string) string {
	num := strings.Trim(number, " \t\r\n")
	num = strings.TrimPrefix(num, "0043")
	num = strings.TrimPrefix(num, "+43")
	num = strings.TrimPrefix(num, "0")
	return num
}
