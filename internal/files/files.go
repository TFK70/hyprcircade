package files

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func ReplaceInFile(filepath string, oldValue string, newValue string, anchor string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	if anchor == "" {
		for i, line := range lines {
			lines[i] = strings.ReplaceAll(line, oldValue, newValue)
		}
	} else {
		for i, line := range lines {
			if strings.Contains(line, anchor) {
				lines[i] = strings.ReplaceAll(line, oldValue, newValue)
			}
		}
	}

	file, err = os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(writer, line)
	}
	writer.Flush()

	return nil
}
