package filefmt

import (
	"archive/zip"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/AuburnCyber/dvsanitizer/internal/safety"
)

func getZipReader(handle *os.File) (*zip.Reader, error) {
	fileInfo, err := handle.Stat()
	if err != nil {
		return nil, fmt.Errorf("could not read attributes: %w", err)
	}
	zipFileSize := fileInfo.Size()

	zipReader, err := zip.NewReader(handle, zipFileSize)
	if err != nil {
		return nil, fmt.Errorf("could not create reader: %w", err)
	}
	// zip.Reader does not need to be closed, only underlying handle.

	return zipReader, nil
}

// Get an ordered-list of all the file-paths in the given zip file.
func ReadInternalZipPaths(zipPath string) []string {
	handle, err := os.Open(zipPath)
	if err != nil {
		safety.ReportError("Error opening zip file.", err,
			"Path: "+zipPath,
		)
	}
	defer handle.Close()

	zipReader, err := getZipReader(handle)
	if err != nil {
		safety.ReportError("Failed to get zip-reader.", err,
			"Path: "+zipPath,
		)
	}

	zipPaths := []string{}
	for _, file := range zipReader.File {
		zipPaths = append(zipPaths, file.Name)
	}

	return zipPaths
}

// Read the given data file from inside the given zip file.
func ReadInternalZipFile(zipPath string, internalPath string) []byte {
	handle, err := os.Open(zipPath)
	if err != nil {
		safety.ReportError("Error opening zip file.", err,
			"Path: "+zipPath,
		)
	}
	defer handle.Close()

	zipReader, err := getZipReader(handle)
	if err != nil {
		safety.ReportError("Failed to get zip-reader.", err,
			"Zip Path: "+zipPath,
		)
	}

	internalHandle, err := zipReader.Open(internalPath)
	if err != nil {
		safety.ReportError("Could not open internal file in JSON-Zip.", err,
			"Zip Path: "+zipPath,
			"Internal Path: "+internalPath,
		)
	}
	defer internalHandle.Close()

	data, err := ioutil.ReadAll(internalHandle)
	if err != nil {
		safety.ReportError("Error encountered reading internal file in JSON-Zip.", err,
			"Zip Path: "+zipPath,
			"Internal Path: "+internalPath,
		)
	}

	return data
}
