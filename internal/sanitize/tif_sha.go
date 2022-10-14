package sanitize

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/AuburnCyber/dvsanitizer/internal/safety"
)

const (
	PATH_TABULATOR_ID_INDEX = 0
	PATH_BATCH_ID_INDEX     = 1
	PATH_RECORD_ID_INDEX    = 2
)

var (
	tifSHAFilenamev1_RE = regexp.MustCompile(`^[0-9]{5}_[0-9]{5}_[0-9]{6}\.(tif|sha)$`)
	tifSHAFilenamev2_RE = regexp.MustCompile(`^[0-9]{5}_[0-9]{5}_[0-9]{6}_[0-9]\.(tif|sha)$`) // For use w/ SF-style filenames.
)

func IsTIFSHAFilename(filename string) bool {
	return tifSHAFilenamev1_RE.Match([]byte(filename)) || tifSHAFilenamev2_RE.Match([]byte(filename))
}

func buildCleanTIFSHAPath(dirtyPath string, dirtyBaseDir string, cleanBaseDir string) string {
	dirtyDir, dirtyFilename := filepath.Split(dirtyPath)
	spFilename := strings.Split(dirtyFilename, ".") // Guaranteed to be 2 via regex enforcement on all only path to reach
	dirtyName, ext := spFilename[0], spFilename[1]
	nameParts := strings.Split(dirtyName, "_")

	// Even though there are 2 formats, the IDs needed are always in the same place
	u32TabID, tabOK := safety.StrToUInt32(nameParts[PATH_TABULATOR_ID_INDEX])
	u32BatchID, batchOK := safety.StrToUInt32(nameParts[PATH_BATCH_ID_INDEX])
	u32RecordID, recordOK := safety.StrToUInt32(nameParts[PATH_RECORD_ID_INDEX])
	if !tabOK || !batchOK || !recordOK {
		safety.ReportError("An incompatible file path was found and no known tactic to recover.", nil,
			"Dirty Path: "+dirtyPath,
			"Tabulator ID: "+nameParts[PATH_TABULATOR_ID_INDEX],
			"Batch ID: "+nameParts[PATH_BATCH_ID_INDEX],
			"Record ID: "+nameParts[PATH_RECORD_ID_INDEX],
		)
	}

	nameParts[2] = createCleanRecordID(u32TabID, u32BatchID, u32RecordID)

	cleanFilename := strings.Join(nameParts, "_") + "." + ext

	outputDirStructure := strings.TrimPrefix(dirtyDir, dirtyBaseDir)

	return filepath.Join(cleanBaseDir, outputDirStructure, cleanFilename)
}

func CleanTIFSHAFiles(dirtyPaths []string, dirtyBaseDir string, cleanBaseDir string) {
	// Create all of the post-sanitization file paths for the .tif/.sha files
	cleanToDirtyDict := make(map[string]string)
	cleanPaths := []string{}
	for _, dirtyPath := range dirtyPaths {
		_, dirtyFilename := filepath.Split(dirtyPath)
		if !IsTIFSHAFilename(dirtyFilename) {
			safety.ReportError("Impermissable filename reached CleanTIFSHAFiles()", nil,
				"Dirty Path: "+dirtyPath,
				"Dirty Filename: "+dirtyFilename,
				"Dirty Base-Dir: "+dirtyBaseDir,
			)
		}

		cleanPath := buildCleanTIFSHAPath(dirtyPath, dirtyBaseDir, cleanBaseDir)

		cleanPaths = append(cleanPaths, cleanPath)
		cleanToDirtyDict[cleanPath] = dirtyPath
	}

	// Unsafe to sanitize the files in-order b/c could leak information about
	// the original filename (which includes the vulnerable record ID) via
	// file-creation timestamp. Instead, process them in a deterministic way
	// based on their new filename which contains the sanitized record ID.
	sort.Strings(cleanPaths)

	for _, cleanPath := range cleanPaths {
		dirtyPath := cleanToDirtyDict[cleanPath]
		dirtyDir, _ := filepath.Split(dirtyPath)
		mode := safety.FileMode(dirtyDir) // Guaranteed to exist b/c walked files to find.
		cleanDir, _ := filepath.Split(cleanPath)

		// Make sure the directory structure exists to write the new .tif/.sha file.
		if err := os.MkdirAll(cleanDir, mode); err != nil {
			safety.ReportError("Unable to build directory structure to write sanitized TIF/SHA file.", err,
				"Clean Path: "+cleanPath,
				"Dirty Path: "+dirtyPath,
				fmt.Sprintf("Mode: %v", mode),
			)
		}

		if err := safety.CopyFile(dirtyPath, cleanPath); err != nil {
			safety.ReportError("Unable to copy TIF/SHA file.", err,
				"Clean Path: "+cleanPath,
				"Dirty Path: "+dirtyPath,
			)
		}
	}

}
