package filefmt

import (
	"errors"
	"fmt"
	"strings"

	"github.com/AuburnCyber/dvsanitizer/internal/safety"
)

func scanField(line string, i int) (string, error) {
	var currField strings.Builder
	isEscaped := false
	for ; i < len(line); i++ {
		switch {
		case line[i] == '"' && isEscaped:
			// Come in pairs only
			isEscaped = false
		case line[i] == '"' && i+1 < len(line) && line[i+1] == '"':
			// Start of an escaped double-quote
			isEscaped = true
		case line[i] == '"':
			// End of field
			return currField.String(), nil
		default:
			// Ensure that doesn't over-escape.
			isEscaped = false
		}

		if err := currField.WriteByte(line[i]); err != nil {
			return "", fmt.Errorf("could not write to string-builder: %w", err)
		}
	}

	return "", errors.New("no field-end indicator")
}

func ConvertCSVLineToFields(line string) ([]string, error) {
	fields := []string{}
	i := 0

	// 1 iteration of the loop for each field in the line
	for i < len(line) {
		// Field may start with `"` or `="`
		switch {
		case strings.HasPrefix(line[i:], `"`):
			i += 1
		case strings.HasPrefix(line[i:], `="`):
			i += 2
		default:
			return nil, fmt.Errorf("no field-prefix at index %d", i)
		}

		if i >= len(line) {
			return nil, errors.New("unexpected end-of-line")
		}

		field, err := scanField(line, i)
		if err != nil {
			safety.StoreErrorDesc("Line: " + line)
			return nil, err
		}

		fields = append(fields, field)

		i += len(field)
		i += 1 // Skip the field-end double-quote

		// Figure-out if there's another field, at the end of the CSV line, or
		// if there's an invalid format.
		switch {
		case i == len(line):
			break // Found end-of-line in expected location
		case line[i] == ',':
			i += 1
			continue // There's another field
		default:
			safety.StoreErrorDesc("Line: " + line)
			return nil, fmt.Errorf("field-separator %q found", line[i])
		}
	}

	return fields, nil
}
