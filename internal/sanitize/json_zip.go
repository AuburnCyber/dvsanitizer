package sanitize

import (
	"archive/zip"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/AuburnCyber/dvsanitizer/internal/filefmt"
	"github.com/AuburnCyber/dvsanitizer/internal/safety"
)

// Files that exist within a JSON-Zip file which can be copied directly b/c
// never seen a record ID in them.
var safeInternalZipFiles = map[string]bool{
	"BallotTypeContestManifest.json":       true,
	"BallotTypeManifest.json":              true,
	"CandidateManifest.json":               true,
	"Configuration.json":                   true,
	"ContestManifest.json":                 true,
	"CountingGroupManifest.json":           true,
	"DistrictManifest.json":                true,
	"DistrictPrecinctPortionManifest.json": true,
	"DistrictTypeManifest.json":            true,
	"ElectionEventManifest.json":           true,
	"OutstackConditionManifest.json":       true,
	"PartyManifest.json":                   true,
	"PrecinctManifest.json":                true,
	"PrecinctPortionManifest.json":         true,
	"TabulatorManifest.json":               true,
}

func createCleanImageMask(dirtyImageMask string, expectedID int, cleanID string) (string, error) {
	// Separate the filename glob from the directory path.
	dir, filenameGlob := filepath.Split(dirtyImageMask) // In the golang source, this walks right-to-left so is safe to use when the filename is a glob
	if len(filenameGlob) == 0 {
		safety.StoreErrorDesc("Dirty Image Mask: " + dirtyImageMask)
		safety.StoreErrorDesc("Dir Portion: " + dir)
		return "", errors.New("unable to separate directory from file-glob in image mask")
	}

	// Split the filename glob the name-elements (before the first star (*)) and after (including the star)
	firstStarIndex := strings.Index(filenameGlob, "*")
	if firstStarIndex == -1 {
		safety.StoreErrorDesc("Dirty Image Mask: " + dirtyImageMask)
		safety.StoreErrorDesc("File Glob: " + filenameGlob)
		return "", errors.New("no star found in file glob")
	}
	beforeStar := filenameGlob[:firstStarIndex]
	afterStar := filenameGlob[firstStarIndex:]

	// Split the pre-star portion of the filename into its individual elements
	nameElements := strings.Split(beforeStar, "_")
	if len(nameElements) < 3 || len(nameElements) > 4 {
		safety.StoreErrorDesc("Dirty Image Mask: " + dirtyImageMask)
		safety.StoreErrorDesc("File Glob: " + filenameGlob)
		safety.StoreErrorDesc(fmt.Sprintf("Split Name: %v", nameElements))
		return "", errors.New("no star found in file glob")
	}

	// Handle the special-case where the glob has an underscore before the star (i.e. '_*.*')
	lastNameElement := nameElements[len(nameElements)-1]
	if len(lastNameElement) == 0 {
		afterStar = "_" + afterStar
		nameElements = nameElements[:len(nameElements)-1]
	}

	// Extract the dirty record ID and consolidate everything before the record ID for later
	dirtyIDStr := nameElements[len(nameElements)-1]
	nameBeforeRecordID := strings.Join(nameElements[:len(nameElements)-1], "_")

	// Make sure that the record ID found in the glob matches the one found in the JSON object
	if dirtyIDStr != fmt.Sprintf("%06d", expectedID) {
		safety.StoreErrorDesc("Found Dirty ID: " + dirtyIDStr)
		safety.StoreErrorDesc(fmt.Sprintf("Expected Dirty ID: %d", expectedID))
		return "", errors.New("non-matching dirty ID found in ImageMask")
	}

	// Reform the image mask using the already calculated clean record ID
	cleanFileGlob := nameBeforeRecordID + "_" + cleanID + afterStar
	return filepath.Join(dir, cleanFileGlob), nil
}

func sanitizeJSONSession(jsonBlob *filefmt.JSONBlob, index int) (string, error) {
	// Pass-through most errors b/c almost nothing to add and only a wrapper to
	// make cleaner due to simpler safety.ReportError() call.

	// Patch-up record ID field.
	tabID, batchID, dirtyRecordID, err := jsonBlob.GetSessionIDs(index)
	if err != nil {
		return "", err
	}

	u32TabID, tabOK := safety.IntToUInt32(tabID)
	u32BatchID, batchOK := safety.IntToUInt32(batchID)
	u32RecordID, recordOK := safety.IntToUInt32(dirtyRecordID)
	if !tabOK || !batchOK || !recordOK {
		safety.StoreErrorDesc("Tabulator ID: " + strconv.Itoa(tabID))
		safety.StoreErrorDesc("Batch ID: " + strconv.Itoa(batchID))
		safety.StoreErrorDesc("Dirty Record ID: " + strconv.Itoa(dirtyRecordID))
		return "", errors.New("incompatible record ID input found")
	}

	cleanRecordID := createCleanRecordID(u32TabID, u32BatchID, u32RecordID)
	if err := jsonBlob.SetSessionRecordID(index, cleanRecordID); err != nil {
		return "", err
	}

	dirtyImageMask, err := jsonBlob.GetSessionImageMask(index)
	if err != nil {
		return "", err
	}

	cleanImageMask, err := createCleanImageMask(dirtyImageMask, dirtyRecordID, cleanRecordID)
	if err != nil {
		safety.StoreErrorDesc("Dirty ImageMask: " + dirtyImageMask)
		return "", err
	}

	if err := jsonBlob.SetSessionImageMask(index, cleanImageMask); err != nil {
		safety.StoreErrorDesc("Dirty ImageMask: " + dirtyImageMask)
		return "", err
	}

	return cleanRecordID, nil
}

func CleanJSONZip(dirtyZipPath string, cleanZipPath string) {
	// Though the pattern is to use the filefmt package to read/write, it's
	// much simpler and lower-complexity to handle writing the zipfile locally.

	if safety.PathExists(cleanZipPath) {
		// This should be caught when parsing command line but saftey-check.
		safety.ReportError("Output file already exists.", nil,
			"Clean Zip Path: "+cleanZipPath,
		)
	}

	cleanZipHandle, err := os.Create(cleanZipPath)
	if err != nil {
		safety.ReportError("Could not create new Zipfile to write to.", err,
			"Clean Zip Path: "+cleanZipPath,
		)
	}
	defer cleanZipHandle.Close()

	cleanZipWriter := zip.NewWriter(cleanZipHandle)
	internalPaths := filefmt.ReadInternalZipPaths(dirtyZipPath)

	// Check to make-sure the provided zip-file is structured as expected to
	// avoid a non-CVR zip being the source file but having a confusing error.
	for _, internalPath := range internalPaths {
		if _, known := safeInternalZipFiles[internalPath]; known {
			continue
		}

		if internalPath == "CvrExport.json" { // Version where there's only 1 file with all the CVRs
			continue
		}

		if strings.HasPrefix(internalPath, "CvrExport_") && strings.HasSuffix(internalPath, ".json") { // Version where there's multiple CVR files
			middlePart := internalPath[len("CvrExport_") : len(internalPath)-len(".json")]
			_, err := strconv.Atoi(middlePart)
			if err == nil {
				continue
			}
		}

		safety.ReportError("The JSON-zip's internal files are not structured as expected.", nil,
			"Unexpected Path: "+internalPath,
			"All Paths: "+strings.Join(internalPaths, " --- "),
		)
	}

	// JSON-Zip files consist of static names (PartyManifest.json,
	// ElectionManifest.json, CvrExport.json, etc) and incrementing names
	// (CvrExport_0.json, CvrExport_1.json, CvrExport_2.json, etc) so don't
	// appear to need sanitization. Process them in alphabetical order anyways
	// for determinism.
	sort.Strings(internalPaths)

	hasSetVersionStr := false
	for _, internalPath := range internalPaths {
		log.Println("PROCESSING INTERNAL ZIPFILE: " + internalPath)

		// Ensure nested directories aren't implicitly included.
		dir, internalFilename := filepath.Split(internalPath)
		if dir != "" || internalFilename != internalPath {
			safety.ReportError(
				"Unexpected zip-internal path", nil,
				"Internal Path: "+internalPath,
				"Dir Portion: "+dir,
				"Filename Portion: "+internalFilename,
			)
		}

		dirtyBytes := filefmt.ReadInternalZipFile(dirtyZipPath, internalFilename)
		var cleanBytes []byte

		if _, isSafe := safeInternalZipFiles[internalFilename]; isSafe {
			log.Println("Skipping safe file")
			cleanBytes = dirtyBytes
		} else {
			// Parse the JSON object
			jsonBlob, err := filefmt.ParseJSON(dirtyBytes)
			if err != nil {
				safety.ReportError("Could not parse JSON extracted from Zip file.", err,
					"Dirty Zip Path: "+dirtyZipPath,
					"Internal Path: "+internalFilename,
				)
			}

			// If this is the first object parsed, note the SW version so it
			// can be included in bug reports.
			if !hasSetVersionStr {
				versionStr := jsonBlob.GetVersionString()
				if versionStr != "" {
					safety.StoreErrorDesc("Software Version: " + versionStr)
					hasSetVersionStr = true
				}
			}

			// Sanitize each voting session in the object and keep track of the
			// clean record IDs created.
			cleanRecordIDs := make([]string, jsonBlob.NumSessions(), jsonBlob.NumSessions())
			for i := 0; i < jsonBlob.NumSessions(); i++ {
				cleanRecordID, err := sanitizeJSONSession(jsonBlob, i)
				if err != nil {
					safety.ReportError("Could not sanitize JSON session.", err,
						"Dirty Zip Path: "+dirtyZipPath,
						"Internal Path: "+internalFilename,
					)
				}
				cleanRecordIDs[i] = cleanRecordID
			}

			// Though not consistent throughout the entire file nor in all
			// files, there appear to be groups of ordered record IDs in the
			// JSON object's Sessions field. To prevent leaking any dirty-ID
			// ordering, reorder the sanitized JSON object's sessions by their
			// clean record ID.
			sort.Strings(cleanRecordIDs)

			cleanBytes, err = filefmt.GetJSONBytes(jsonBlob, cleanRecordIDs)
			if err != nil {
				safety.ReportError("Could not marshal sanitized JSON object.", err,
					"Internal Path: "+internalFilename,
					"Dirty Zip Path: "+dirtyZipPath,
				)
			}
		}

		internalWriter, err := cleanZipWriter.Create(internalFilename)
		if err != nil {
			safety.ReportError("Could not create writer to store data inside cleaned JSON zip.", err,
				"Clean Zip Path: "+cleanZipPath,
				"Internal filename: "+internalFilename,
			)
		}
		// Internal writer doesn't have a Close() function b/c handled by the zip-writer's Close()

		if n, err := internalWriter.Write(cleanBytes); n != len(cleanBytes) || err != nil {
			safety.ReportError("Failed to write clean JSON bytes to clean zip file.", err,
				"Clean Zip Path: "+cleanZipPath,
				"Internal filename: "+internalFilename,
				"n: "+strconv.Itoa(n),
			)
		}
	}

	if err := cleanZipWriter.Close(); err != nil {
		safety.ReportError("Unable to close clean Zipfile writer.", err,
			"Clean Zip Path: "+cleanZipPath,
		)
	}

}
