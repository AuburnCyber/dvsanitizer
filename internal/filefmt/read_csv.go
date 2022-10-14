package filefmt

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/hex"
	"fmt"
	"io"
	"unicode/utf8"

	"github.com/AuburnCyber/dvsanitizer/internal/safety"
)

// Turn the CSV file's byte into a matrix of fields of the string-value // contents.
func ReadCSVBytes(csvBytes []byte) ([][]string, bool, error) {
	var lines [][]string
	var isOddballFormat bool
	var err error

	if csvBytes[0] == '"' {
		isOddballFormat = true
		lines, err = readOddballCSVFile(csvBytes)
	} else { // Will over-catch invalid formats but should be handled by the golang CSV library.
		isOddballFormat = false
		lines, err = readStandardCSVFile(csvBytes)
	}

	if err == nil {
		baseWidth := len(lines[0])
		for i := 0; i < len(lines); i++ {
			if len(lines[i]) != baseWidth {
				safety.StoreErrorDesc(fmt.Sprintf("First Line: %v", lines[0]))
				safety.StoreErrorDesc(fmt.Sprintf("Diff Line: %v", lines[i]))

				return nil, false, fmt.Errorf("non-square CSV read. first-line: %d, diff-line: %d, diff-line-index: %d", baseWidth, len(lines[i]), i)
			}
		}
	}

	return lines, isOddballFormat, err
}

func readOddballCSVFile(csvBytes []byte) ([][]string, error) {
	readLines := [][]string{}

	scanner := bufio.NewScanner(bytes.NewBuffer(csvBytes))
	for scanner.Scan() {
		inLine := scanner.Text()
		if !utf8.Valid([]byte(inLine)) {
			safety.StoreErrorDesc("Invalid Unicode Line (hex-encoded): " + hex.EncodeToString([]byte(inLine)))
			return nil, fmt.Errorf("invalid unicode line in input. line-index: %d", len(readLines))
		}

		fields, err := ConvertCSVLineToFields(scanner.Text())
		if err != nil {
			safety.StoreErrorDesc("Line: " + scanner.Text())
			return nil, fmt.Errorf("CSV line-conversion error, line-index: %d, err: %w", len(readLines), err)
		}
		readLines = append(readLines, fields)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("CSV scan error: %w", err)
	}

	return readLines, nil
}

func readStandardCSVFile(csvBytes []byte) ([][]string, error) {
	readLines := [][]string{}
	reader := csv.NewReader(bytes.NewReader(csvBytes))

	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil { // other than EOF
			return nil, fmt.Errorf("failed to read STD CSV line: %w", err)
		}
		readLines = append(readLines, line)

	}

	return readLines, nil
}
