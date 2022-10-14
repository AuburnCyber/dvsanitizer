package filefmt

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/AuburnCyber/dvsanitizer/internal/safety"
)

// A type to abstract away the awfulness that is partial-decoding of the
// complete JSON object.
type JSONBlob struct {
	root map[string]json.RawMessage

	// These fields are index-aligned to prevent unmarshalling every session
	// multiple times when re-ordering sessions.
	originalOrderSessionList []json.RawMessage
	originalOrderCleanIDs    []string
}

// The number of JSON objects in ".Sessions" of the root object.
func (b *JSONBlob) NumSessions() int {
	return len(b.originalOrderSessionList)
}

// Extract the tabulator ID, batch ID, and dirty record ID of the indexed
// Session object for sanitization.
func (b *JSONBlob) GetSessionIDs(index int) (retTabID int, retBatchID int, retRecordID int, retErr error) {
	if index >= b.NumSessions() {
		safety.StoreErrorDesc("Non-Existent Index: " + strconv.Itoa(index))
		return -1, -1, -1, errors.New("out-of-bounds session index")
	}
	sessionRaw := b.originalOrderSessionList[index]

	var session map[string]json.RawMessage
	if err := json.Unmarshal(sessionRaw, &session); err != nil {
		safety.StoreErrorDesc("Index: " + strconv.Itoa(index))
		return -1, -1, -1, fmt.Errorf("unable to unmarshal session object: %w", err)
	}

	var tabID int
	if err := json.Unmarshal(session["TabulatorId"], &tabID); err != nil {
		safety.StoreErrorDesc("Index: " + strconv.Itoa(index))
		safety.StoreErrorDesc("Session Object Causing Error (base64 encoded bytes): " + base64.StdEncoding.EncodeToString(sessionRaw))
		return -1, -1, -1, fmt.Errorf("unable to unmarshal tabulator ID: %w", err)
	}

	var batchID int
	if err := json.Unmarshal(session["BatchId"], &batchID); err != nil {
		safety.StoreErrorDesc("Index: " + strconv.Itoa(index))
		safety.StoreErrorDesc("Session Object Causing Error (base64 encoded bytes): " + base64.StdEncoding.EncodeToString(sessionRaw))
		return -1, -1, -1, fmt.Errorf("unable to unmarshal batch ID: %w", err)
	}

	var recordID int
	if err := json.Unmarshal(session["RecordId"], &recordID); err != nil {
		safety.StoreErrorDesc("Index: " + strconv.Itoa(index))
		safety.StoreErrorDesc("Session Object Causing Error (base64 encoded bytes): " + base64.StdEncoding.EncodeToString(sessionRaw))
		return -1, -1, -1, fmt.Errorf("unable to unmarshal recordID: %w", err)
	}

	return tabID, batchID, recordID, nil
}

// Replace the dirty record ID of the indexed Session object with the provided
// post-sanitization record ID.
// func (b *JSONBlob) SetSessionRecordID(index int, newRecordID int) error {
func (b *JSONBlob) SetSessionRecordID(index int, newRecordID string) error {
	if index >= b.NumSessions() {
		safety.StoreErrorDesc("Non-Existent Index: " + strconv.Itoa(index))
		return errors.New("out-of-bounds session index")
	}

	recordIDRaw, err := json.Marshal(newRecordID)
	if err != nil {
		safety.StoreErrorDesc("Clean Record ID: " + newRecordID)
		return fmt.Errorf("unable to marshal raw version of clean record ID: %w", err)
	}

	dirtySessionRaw := b.originalOrderSessionList[index]

	var session map[string]json.RawMessage
	if err := json.Unmarshal(dirtySessionRaw, &session); err != nil {
		safety.StoreErrorDesc("Index: " + strconv.Itoa(index))
		return fmt.Errorf("unable to unmarshal session object: %w", err)
	}

	session["RecordId"] = recordIDRaw

	cleanSessionRaw, err := json.Marshal(session)
	if err != nil {
		safety.StoreErrorDesc("Clean Record ID: " + newRecordID)
		safety.StoreErrorDesc("Clean Record ID Object (base64 encoded bytes): " + base64.StdEncoding.EncodeToString(recordIDRaw))
		safety.StoreErrorDesc("Session Object (base64 encoded bytes): " + base64.StdEncoding.EncodeToString(dirtySessionRaw))
		return fmt.Errorf("unable to marshal session object with clean record ID: %w", err)
	}

	b.originalOrderSessionList[index] = cleanSessionRaw
	b.originalOrderCleanIDs[index] = newRecordID

	return nil
}

// Extract the dirty ImageMask from the indexed Session object for sanitization.
func (b *JSONBlob) GetSessionImageMask(index int) (string, error) {
	if index >= b.NumSessions() {
		safety.StoreErrorDesc("Non-Existent Index: " + strconv.Itoa(index))
		return "", errors.New("out-of-bounds session index")
	}

	sessionRaw := b.originalOrderSessionList[index]

	var session map[string]json.RawMessage
	if err := json.Unmarshal(sessionRaw, &session); err != nil {
		safety.StoreErrorDesc("Index: " + strconv.Itoa(index))
		safety.StoreErrorDesc("Session Object Causing Error (base64 encoded bytes): " + base64.StdEncoding.EncodeToString(sessionRaw))
		return "", fmt.Errorf("could not unmarshal raw session object: %w", err)
	}

	var imageMask string
	if err := json.Unmarshal(session["ImageMask"], &imageMask); err != nil {
		safety.StoreErrorDesc("Index: " + strconv.Itoa(index))
		safety.StoreErrorDesc("Session Object Causing Error (base64 encoded bytes): " + base64.StdEncoding.EncodeToString(sessionRaw))
		return "", fmt.Errorf("could not unmarshal session object's ImageMask: %v", err)
	}

	return imageMask, nil
}

// Replace the dirty ImageMask of the indexed Session object with the provided
// post-sanitization ImageMask.
func (b *JSONBlob) SetSessionImageMask(index int, cleanImageMask string) error {
	if index >= b.NumSessions() {
		safety.StoreErrorDesc("Non-Existent Index: " + strconv.Itoa(index))
		return errors.New("out-of-bounds session index")
	}

	dirtySessionRaw := b.originalOrderSessionList[index]

	imageMaskRaw, err := json.Marshal(cleanImageMask)
	if err != nil {
		safety.StoreErrorDesc("Clean ImageMask: " + cleanImageMask)
		return fmt.Errorf("unable to marshal raw version of clean record ID: %w", err)
	}

	var session map[string]json.RawMessage
	if err := json.Unmarshal(dirtySessionRaw, &session); err != nil {
		safety.StoreErrorDesc("Index: " + strconv.Itoa(index))
		safety.StoreErrorDesc("Clean ImageMask: " + cleanImageMask)
		return fmt.Errorf("unable to unmarshal dirty session: %w", err)
	}

	session["ImageMask"] = imageMaskRaw

	sessionRaw, err := json.Marshal(session)
	if err != nil {
		safety.StoreErrorDesc("Index: " + strconv.Itoa(index))
		safety.StoreErrorDesc("Clean ImageMask: " + cleanImageMask)
		return fmt.Errorf("unable to marshal clean session: %w", err)
	}

	b.originalOrderSessionList[index] = sessionRaw

	return nil
}

// Get the version string from the root-level of the object to assist in remote
// debugging. Unlike other functions, this one is allowed to fail for any
// reason and will return "NOT-AVAILABLE".
func (b *JSONBlob) GetVersionString() string {
	var versionStr string
	if err := json.Unmarshal(b.root["Version"], &versionStr); err != nil {
		log.Printf("unable to read version string from JSON")
		return "NOT-AVAILABLE"
	}

	return versionStr
}

// Panic because it should never be called.
func (b *JSONBlob) MarshalJSON() ([]byte, error) {
	// Is a coding error but don't let it hide if happens again.
	panic("JSONBlob.MarshalJSON() not implemented, use GetJSONBytes())")
}

// Convert the JSON bytes into a form that can be easily handled for
// sanitization of the internal JSON objects.
func ParseJSON(jsonBytes []byte) (*JSONBlob, error) {
	blob := &JSONBlob{}

	if err := json.Unmarshal(jsonBytes, &blob.root); err != nil {
		return nil, fmt.Errorf("root-unmarshal error: %w", err)
	}

	_, ok := blob.root["Sessions"]
	if !ok {
		return nil, errors.New("JSON object does not have a 'Sessions' key")
	}

	if err := json.Unmarshal(blob.root["Sessions"], &blob.originalOrderSessionList); err != nil {
		return nil, fmt.Errorf("sessions-unmarshal error: %w", err)
	}

	// Used to greatly simplify re-ordering when marshalling
	blob.originalOrderCleanIDs = make([]string, len(blob.originalOrderSessionList), len(blob.originalOrderSessionList))

	return blob, nil
}

// Unmarshalling all the session objects multiple times to fetch their clean
// record IDs is painful so use the post-sanitization record IDs that were
// saved when the objects were updated to re-order the marshalled session bytes
// without having to unmarshal.
func reorderSessionList(preOrderingSessions []json.RawMessage, preOrderingRecordIDs []string, newOrderingRecordIDs []string) []json.RawMessage {
	reorderedSessions := make([]json.RawMessage, len(preOrderingSessions), len(preOrderingSessions))

	for postOrderingIndex, nextRecordID := range newOrderingRecordIDs {
		nextIndexToAdd := -1
		for preOrderingIndex := 0; preOrderingIndex < len(preOrderingSessions); preOrderingIndex++ {
			if preOrderingRecordIDs[preOrderingIndex] == nextRecordID {
				nextIndexToAdd = preOrderingIndex
				break
			}
		}
		if nextIndexToAdd == -1 {
			safety.ReportError("Could not re-order JSON sessions.", nil,
				"Unfound ID: "+nextRecordID,
				"Cleaned Record IDs: "+strings.Join(preOrderingRecordIDs, ", "),
				"New Order: "+strings.Join(newOrderingRecordIDs, ", "),
			)
		}

		reorderedSessions[postOrderingIndex] = preOrderingSessions[nextIndexToAdd]
	}

	return reorderedSessions
}

// Convert the object to JSON bytes after re-ordering the Sessions list based on recordIDOrder.
func GetJSONBytes(jsonBlob *JSONBlob, recordIDOrder []string) ([]byte, error) {
	reorderedSessions := reorderSessionList(jsonBlob.originalOrderSessionList, jsonBlob.originalOrderCleanIDs, recordIDOrder)

	rawSessions, err := json.Marshal(reorderedSessions)
	if err != nil {
		return nil, fmt.Errorf("session-list marshal error: %w", err)
	}

	jsonBlob.root["Sessions"] = rawSessions

	rawRoot, err := json.Marshal(jsonBlob.root)
	if err != nil {
		return nil, fmt.Errorf("root marshal error: %w", err)
	}

	return []byte(rawRoot), nil
}
