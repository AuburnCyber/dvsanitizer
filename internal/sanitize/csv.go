package sanitize

import (
	"errors"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/AuburnCyber/dvsanitizer/internal/filefmt"
	"github.com/AuburnCyber/dvsanitizer/internal/safety"
)

const (
	CSV_NUM_HEADER_ROWS = 4
	CSV_DATA_ROWS_START = CSV_NUM_HEADER_ROWS

	CSV_CVR_NUM_INDEX      = 0
	CSV_TABULATOR_ID_INDEX = 1
	CSV_BATCH_ID_INDEX     = 2
	CSV_RECORD_ID_INDEX    = 3
)

func resortCSVDataLines(dataLines [][]string) [][]string {
	// Re-order the data-lines based on tabulator + batch + new record ID.
	sort.Slice(dataLines, func(leftIndex int, rightIndex int) bool {
		left := dataLines[leftIndex]
		right := dataLines[rightIndex]

		// Will pad with spaces (0x20) so is still ordered correctly when compared
		leftSortable := fmt.Sprintf("%15s%15s%15s", left[CSV_TABULATOR_ID_INDEX], left[CSV_BATCH_ID_INDEX], left[CSV_RECORD_ID_INDEX])
		rightSortable := fmt.Sprintf("%15s%15s%15s", right[CSV_TABULATOR_ID_INDEX], right[CSV_BATCH_ID_INDEX], right[CSV_RECORD_ID_INDEX])

		return leftSortable < rightSortable
	})

	// Remove the CVR Numbers since they no-longer map to the same record ID.
	for i := 0; i < len(dataLines); i++ {

		dataLines[i][0] = ""
	}

	return dataLines
}

func validateReadCSV(dirtyLines [][]string) error {
	if len(dirtyLines) >= 1 && len(dirtyLines[0]) >= 2 {
		// Though save so that is included in any error-report going forward.
		safety.StoreErrorDesc("Software Version: " + dirtyLines[0][1])
	}
	// Check the attributes of the read CSV to make sure fits expectations.
	switch {
	case len(dirtyLines) < 4:
		safety.StoreErrorDesc("Num Rows: " + strconv.Itoa(len(dirtyLines)))
		return errors.New("CSV is too small for expected header rows")

	case dirtyLines[3][CSV_CVR_NUM_INDEX] != "CvrNumber":
		safety.StoreErrorDesc("Found header: " + dirtyLines[3][CSV_CVR_NUM_INDEX])
		return errors.New("CvrNumber row-header incorrect")

	case dirtyLines[3][CSV_TABULATOR_ID_INDEX] != "TabulatorNum":
		safety.StoreErrorDesc("Found header: " + dirtyLines[3][CSV_TABULATOR_ID_INDEX])
		return errors.New("TabulatorNum row-header incorrect")

	case dirtyLines[3][CSV_BATCH_ID_INDEX] != "BatchId":
		safety.StoreErrorDesc("Found header: " + dirtyLines[3][CSV_BATCH_ID_INDEX])
		return errors.New("BatchId row-header incorrect")

	case dirtyLines[3][CSV_RECORD_ID_INDEX] != "RecordId":
		safety.StoreErrorDesc("Found header: " + dirtyLines[3][CSV_RECORD_ID_INDEX])
		return errors.New("RecordId row-header incorrect")

	case len(dirtyLines) < CSV_NUM_HEADER_ROWS+1:
		safety.StoreErrorDesc("Num Rows: " + strconv.Itoa(len(dirtyLines)))
		return errors.New("CSV is too small for any data rows")
	}

	return nil
}

// Read the dirtyPath CSV in either standard CSV or odd-ball CSV (="0" style fields) format and then:
//		- Sanitize the record ID column
//		- Re-order rows based on the new record ID
//			- Tabulator and batch IDs are maintaine
//		- Re-number the CVR Numbers of the newly-ordered rows
//		- Write to cleanPath with the same format that was read (standard or odd-ball)
func CleanCSVFile(dirtyPath string, cleanPath string) {
	if !safety.IsFile(dirtyPath) {
		safety.ReportError("Attempting to clean a non-existant CSV file.", nil,
			"Dirty Path: "+dirtyPath,
		)
	}
	dirtyBytes, err := os.ReadFile(dirtyPath)
	if err != nil {
		safety.ReportError("Unable to read dirty CSV file.", err,
			"Dirty Path: "+dirtyPath,
		)
	}

	allDirtyLines, isOddballFormat, err := filefmt.ReadCSVBytes(dirtyBytes)
	if err != nil {
		safety.ReportError("Unable to parse CSV file.", err,
			"Dirty Path: "+dirtyPath,
		)
	}
	if isOddballFormat {
		log.Printf("Found %d rows of data in oddball CSV format.", len(allDirtyLines))
	} else {
		log.Printf("Found %d rows of data in standard CSV format.", len(allDirtyLines))
	}

	if err := validateReadCSV(allDirtyLines); err != nil {
		safety.ReportError("Invalid or unhandled CSV construction can not be sanitized.", nil,
			"Dirty Path: "+dirtyPath,
		)
	}

	headerLines := make([][]string, CSV_NUM_HEADER_ROWS, CSV_NUM_HEADER_ROWS)
	copy(headerLines, allDirtyLines)

	dirtyLines := make([][]string, len(allDirtyLines)-CSV_NUM_HEADER_ROWS, len(allDirtyLines)-CSV_NUM_HEADER_ROWS)
	copy(dirtyLines, allDirtyLines[CSV_NUM_HEADER_ROWS:])

	cleanLines, err := cleanCSVLines(dirtyLines)
	if err != nil {
		safety.ReportError("Could not sanitize CSV lines.", err,
			"Dirty Path: "+dirtyPath,
		)
	}

	resortedCleanLines := resortCSVDataLines(cleanLines)

	allCleanLines := append(headerLines, resortedCleanLines...)

	cleanBytes, err := filefmt.LinesToCSVBytes(allCleanLines, isOddballFormat)
	if err != nil {
		safety.ReportError("Unable to marshal CSV bytes.", err,
			"Dirty Path: "+dirtyPath,
		)
	}

	if err := os.WriteFile(cleanPath, cleanBytes, safety.FileMode(dirtyPath)); err != nil {
		safety.ReportError("Unable to write cleaned CSV to file.", err,
			"Dirty Path: "+dirtyPath,
			"Clean Path: "+cleanPath,
		)
	}
}

func cleanCSVLines(dirtyLines [][]string) ([][]string, error) {
	cleanLines := make([][]string, len(dirtyLines), len(dirtyLines))
	for i := 0; i < len(dirtyLines); i++ {
		line := dirtyLines[i]

		u32TabID, tabOK := safety.StrToUInt32(line[CSV_TABULATOR_ID_INDEX])
		u32BatchID, batchOK := safety.StrToUInt32(line[CSV_BATCH_ID_INDEX])
		u32RecordID, recordOK := safety.StrToUInt32(line[CSV_RECORD_ID_INDEX])
		if !tabOK || !batchOK || !recordOK {
			safety.StoreErrorDesc("Dirty Line: " + strings.Join(line, " --- "))
			safety.StoreErrorDesc("Tabulator ID: " + line[CSV_TABULATOR_ID_INDEX])
			safety.StoreErrorDesc("Batch ID: " + line[CSV_BATCH_ID_INDEX])
			safety.StoreErrorDesc("Record ID: " + line[CSV_RECORD_ID_INDEX])
			return nil, errors.New("incompatible record ID input found")
		}
		cleanRecordID := createCleanRecordID(u32TabID, u32BatchID, u32RecordID)
		line[CSV_RECORD_ID_INDEX] = cleanRecordID
		cleanLines[i] = line
	}

	return cleanLines, nil
}
