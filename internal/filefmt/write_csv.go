package filefmt

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"strings"
	"unicode"

	"github.com/AuburnCyber/dvsanitizer/internal/safety"
)

// Attempts to mimic the behavior of the EMS's CSV writer but is not guaranteed to be 100% correct.
func LinesToCSVBytes(cleanLines [][]string, isOddballFormat bool) ([]byte, error) {
	baseWidth := len(cleanLines[0])
	for i := 0; i < len(cleanLines); i++ {
		if len(cleanLines[i]) != baseWidth {
			safety.StoreErrorDesc(fmt.Sprintf("First Line: %v", cleanLines[0]))
			safety.StoreErrorDesc(fmt.Sprintf("Diff Line: %v", cleanLines[i]))

			return nil, fmt.Errorf("non-square CSV created. first-line: %d, diff-line: %d, diff-line-index: %d", baseWidth, len(cleanLines[i]), i)
		}
	}

	if isOddballFormat {
		return linesToOddballCSVBytes(cleanLines)
	} else {
		return linesToStandardCSVBytes(cleanLines)
	}
}

func fieldWantsPrefix(field string) bool {
	if len(field) == 0 {
		return false
	}

	for _, aRune := range field {
		if aRune != '-' && !unicode.IsDigit(aRune) {
			return false
		}
	}
	return true
}

func linesToOddballCSVBytes(cleanLines [][]string) ([]byte, error) {
	var buffer bytes.Buffer

	for i, line := range cleanLines {
		if i < 4 {
			// Header-lines are encoded differently from data-lines.
			marshalledLine := `"`
			marshalledLine += strings.Join(line, `","`)
			marshalledLine += `"` + "\r\n"
			if n, err := buffer.WriteString(marshalledLine); n != len(marshalledLine) || err != nil {
				safety.StoreErrorDesc(fmt.Sprintf("Header Line: %v", line))
				return nil, fmt.Errorf("could not write header-line to buffer. n: %d, err:%w", n, err)
			}
			continue
		}

		for fieldIndex, field := range line {
			var out string
			switch {
			case fieldIndex < 5 && fieldWantsPrefix(field):
				out = `="` + field + `"`
			case fieldIndex < 5 && field == "":
				out = `=""`
			case len(field) == 0:
				out = `""`
			default:
				out = `"` + field + `"`
			}
			if n, err := buffer.WriteString(out); n != len(out) || err != nil {
				return nil, fmt.Errorf("unable to write field to buffer. field: %q, n: %d, err: %v", field, n, err)
			}

			if fieldIndex != len(line)-1 {
				if err := buffer.WriteByte(','); err != nil {
					return nil, fmt.Errorf("unable to write field-separator to buffer: %v", err)
				}
			}

		}
		if n, err := buffer.WriteString("\r\n"); n != 2 || err != nil {
			return nil, fmt.Errorf("unable to write newline to buffer. n: %d, err: %v", n, err)
		}
	}

	return buffer.Bytes(), nil
}

func linesToStandardCSVBytes(cleanLines [][]string) ([]byte, error) {
	var buffer bytes.Buffer
	writer := csv.NewWriter(&buffer)

	for _, line := range cleanLines {
		if err := writer.Write(line); err != nil {
			safety.StoreErrorDesc(fmt.Sprintf("Line: %v", line))
			return nil, fmt.Errorf("could not write to CSV writer: %w", err)
		}
	}

	writer.Flush() // Doesn't return error

	return buffer.Bytes(), nil

}
