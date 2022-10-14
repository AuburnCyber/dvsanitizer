package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/AuburnCyber/dvsanitizer/internal/sanitize"
)

const (
	UNFINISHED_DIR_PREFIX = "SANITIZATION-UNFINISHED_"

	GEN_SEED_SIZE = 16 // in bytes

	MSG_BREAK = `********************************************************************************
********************************************************************************`

	GEN_SEED_MSG = "" +
		`Sanitization succeeded with the above auto-generated seed.

If you wish to sanitize multiple files and keep sanitized record IDs
consistent, you MUST provide the below auto-generated seed to the next instance
by using the '--seed' flag. Reasons you might want to do this are:
    You want to sanitize both the CSV and JSON-zip format CVR files and be able
    to reference a single ballot with a single sanitized record ID.
        First call : --sanitize-csv         --gen-seed
        Second call: --sanitize-json-zip    --seed=XXXXXXXXXXXXXXXX

    You want to sanitze a CSV CVR file and its associated ballot images.
        First instance : --sanitize-csv     --gen-seed
        Second instance: --sanitize-tif-dir --seed=XXXXXXXXXXXXXXXX

    You want to sanitize a partial CVR file then a full CVR file.
        First instance : --sanitize-csv     --gen-seed
        Second instance: --sanitize-csv     --seed=XXXXXXXXXXXXXXXX`
)

func createUnfinishedDir() string {
	unfinishedDirPath := filepath.Join(
		args.outputDir,
		UNFINISHED_DIR_PREFIX+time.Now().Format("2006-01-02T15-04-05"),
	)

	if err := os.Mkdir(unfinishedDirPath, 0755); err != nil {
		log.Fatalf("Unable to create directory for in-progress sanitization: %v", err)
	}

	return unfinishedDirPath
}

func main() {
	ParseCmd()

	if args.genSeed {
		rawBuffer := make([]byte, GEN_SEED_SIZE, GEN_SEED_SIZE)
		if n, err := rand.Read(rawBuffer); n != GEN_SEED_SIZE || err != nil {
			log.Fatalf("An error was encountered generating a secure random seed.")
		}
		args.seed = hex.EncodeToString(rawBuffer)
	}
	sanitize.Initialize([]byte(args.seed))

	var unfinishedDirPath string
	switch {
	case args.actionSanitizeCSV:
		validateSanitizeCSVFlags()
		unfinishedDirPath = createUnfinishedDir()
		SanitizeCSV(unfinishedDirPath)
	case args.actionSanitizeJSONZip:
		validateSanitizeJSONZipFlags()
		unfinishedDirPath = createUnfinishedDir()
		SanitizeJSONZip(unfinishedDirPath)
	case args.actionSanitizeTIFDir:
		validateSanitizeTIFDirFlags()
		unfinishedDirPath = createUnfinishedDir()
		SanitizeTIFDir(unfinishedDirPath)
	}

	// Now that we're finished, move everything out of the unfinished directory
	// into the actual output directory.
	entries, err := os.ReadDir(unfinishedDirPath)
	if err != nil {
		log.Fatalf("Unable to find contents of the unfinished-directory containing the post-sanitization files: %v", err)
	}
	for _, entry := range entries {
		unfinishedPath := filepath.Join(unfinishedDirPath, entry.Name())
		finishedPath := filepath.Join(args.outputDir, entry.Name())

		if err := os.Rename(unfinishedPath, finishedPath); err != nil {
			log.Fatalf("Unable to mark file as finished. unfinished-path: %q, finished-path: %q", unfinishedPath, finishedPath)
		}
	}
	if err := os.Remove(unfinishedDirPath); err != nil {
		log.Fatalf("Unable to delete now-empty unfinished-path: %v", err)
	}

	if args.genSeed {
		fmt.Println(MSG_BREAK)
		fmt.Println("Your auto-generated seed is: " + args.seed)
		fmt.Println(MSG_BREAK)
		fmt.Println(GEN_SEED_MSG)
		fmt.Println(MSG_BREAK)
	}
}
