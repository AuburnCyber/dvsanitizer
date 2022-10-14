package filefmt

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/AuburnCyber/dvsanitizer/internal/safety"
)

func WriteSkippedFiles(outputPath string, wrongExtPaths []string, wrongNameFormatPaths []string, unkSafetyPaths []string) error {
	if safety.PathExists(outputPath) {
		// Checked in tif_actions.go but is easy to make sure.
		panic("skipped-file already exists")
	}

	handle, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("unable to create skipped-file: %v", err)
	}
	defer handle.Close()
	writer := csv.NewWriter(handle)

	if err := writer.Write([]string{"Skip Reason", "Path"}); err != nil {
		return fmt.Errorf("unable to write header: %v", err)
	}

	for _, path := range unkSafetyPaths {
		if err := writer.Write([]string{"Unknown safety of contents", path}); err != nil {
			return fmt.Errorf("unable to write unk-safety line: %v", err)
		}
	}
	for _, path := range wrongNameFormatPaths {
		if err := writer.Write([]string{"Unrecognized TIF/SHA filename format", path}); err != nil {
			return fmt.Errorf("unable to write unk-format line: %v", err)
		}
	}
	for _, path := range wrongExtPaths {
		if err := writer.Write([]string{"Not a .tif or .sha file", path}); err != nil {
			return fmt.Errorf("unable to write wrong-extension line: %v", err)
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return fmt.Errorf("post-flush error encountered: %v", err)
	}

	return nil
}
