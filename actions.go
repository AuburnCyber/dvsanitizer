package main

import (
	"log"
	"path/filepath"
	"strings"

	"github.com/AuburnCyber/dvsanitizer/internal/filefmt"
	"github.com/AuburnCyber/dvsanitizer/internal/safety"
	"github.com/AuburnCyber/dvsanitizer/internal/sanitize"
)

func SanitizeTIFDir(unfinishedDirPath string) {
	// Get a list of all the files and directories in the input directory.
	allDirtyPaths, err := safety.GetFileListing(args.inputPath)
	if err != nil {
		safety.ReportError("Unable get listing of files in TIF-dir.", err,
			"Input Path: "+args.inputPath,
		)
	}

	// Filter down to only the files that need to be sanitized.
	toBeSanitizedPaths := []string{}
	var wrongExtPaths []string
	var wrongNameFormatPaths []string
	var unkSafetyPaths []string
	for _, dirtyPath := range allDirtyPaths {
		_, dirtyFilename := filepath.Split(dirtyPath)

		switch {
		case filepath.Ext(dirtyFilename) != ".tif" && filepath.Ext(dirtyFilename) != ".sha":
			wrongExtPaths = append(wrongExtPaths, dirtyPath)

		case !sanitize.IsTIFSHAFilename(dirtyFilename):
			wrongNameFormatPaths = append(wrongNameFormatPaths, dirtyPath)

		case strings.HasPrefix(dirtyFilename, "NotCast") || strings.HasPrefix(dirtyFilename, "_NotCast"):
			// Seen a handful of times in the GA data and have not been able to
			// confirm that they do not contain record IDs or other metadata
			// within the TIF structure.
			unkSafetyPaths = append(unkSafetyPaths, dirtyPath)

		default:
			toBeSanitizedPaths = append(toBeSanitizedPaths, dirtyPath)

		}
	}

	skippedPath := filepath.Join(unfinishedDirPath, "skipped-during-sanitization.csv")
	if safety.PathExists(skippedPath) {
		log.Fatalf("Refusing to execute as skipped-file already exists: %s", skippedPath)
	}
	if err := filefmt.WriteSkippedFiles(skippedPath, wrongExtPaths, wrongNameFormatPaths, unkSafetyPaths); err != nil {
		log.Fatalf("Encountered error writing to skipped-file: %v", err)
	}

	sanitize.CleanTIFSHAFiles(toBeSanitizedPaths, args.inputPath, unfinishedDirPath)
}

func SanitizeJSONZip(unfinishedDirPath string) {
	_, filename := filepath.Split(args.inputPath)
	cleanZipPath := filepath.Join(unfinishedDirPath, filename)
	sanitize.CleanJSONZip(args.inputPath, cleanZipPath)
}

func SanitizeCSV(unfinishedDirPath string) {
	_, filename := filepath.Split(args.inputPath)
	cleanCSVPath := filepath.Join(unfinishedDirPath, filename)
	sanitize.CleanCSVFile(args.inputPath, cleanCSVPath)
}
