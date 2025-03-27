package utils

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"sync"
)

var csvMutex sync.Mutex

func FileReader(filename string) ([]string, error) {
	var lines []string
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}

func CsvWriter(headers, values []string, filePath string) error {
	csvMutex.Lock()
	defer csvMutex.Unlock()

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)

	info, err := file.Stat()
	if err != nil {
		return err
	}

	if info.Size() == 0 && len(headers) > 0 {
		if err := writer.Write(headers); err != nil {
			return fmt.Errorf("failed to write header: %w", err)
		}
	}

	if err := writer.Write(values); err != nil {
		return fmt.Errorf("failed to write row: %w", err)
	}

	writer.Flush()
	return writer.Error()
}
