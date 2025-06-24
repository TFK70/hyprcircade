package files

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/tfk70/hyprcircade/internal/logging"
)

func ReplaceInFile(filepath string, oldValue string, newValue string, anchor string) error {
	logger, err := logging.GetNamedLogger("files.go")
	if err != nil {
		return err
	}

	fileLogger := logger.WithFields(logrus.Fields{
		"path":     filepath,
		"oldValue": oldValue,
		"newValue": newValue,
		"anchor":   anchor,
	})

	fileLogger.Info("Performing replacement")

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

	fileLogger.Info("Successfully replaced")

	return nil
}
