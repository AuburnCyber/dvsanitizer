package filefmt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func Test_JSONBlob_NoSessions(t *testing.T) {
	testcases := []string{
		`{"Version":"1.2.3","Sessions":null}`,
		`{"Version":"1.2.3","Sessions":[]}`,
	}

	for i, tc := range testcases {
		t.Run(fmt.Sprintf("Index-%d", i), func(t *testing.T) {
			got, err := ParseJSON([]byte(tc))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			_, _, _, session0Err := got.GetSessionIDs(0)
			_, image0Err := got.GetSessionImageMask(0)
			_, _, _, session1Err := got.GetSessionIDs(1)
			_, image1Err := got.GetSessionImageMask(1)

			switch {
			case got.NumSessions() != 0:
				t.Errorf("incorrect number of sessions")
			case got.GetVersionString() != "1.2.3":
				t.Errorf("incorrect version string")
			case session0Err == nil:
				t.Errorf("incorrect session 0 fetch")
			case image0Err == nil:
				t.Errorf("incorrect image-mask 0 fetch")
			case session1Err == nil:
				t.Errorf("incorrect session 1 fetch")
			case image1Err == nil:
				t.Errorf("incorrect image-mask 1 fetch")
			}
		})
	}
}

func Test_JSONBlob_SessionIDs(t *testing.T) {
	// This is too much under-test in a single testcase but whole lot of
	// boiler-plat to break-up
	input := []byte(
		`{"Version":"4.5.6","Sessions":[{` +
			`"TabulatorId":1,` +
			`"BatchId":2,` +
			`"RecordId":3,` +
			`"ImageMask":"ABC"` +
			`},{` +
			`"Stuff-1":"Begin",` +
			`"TabulatorId":4,` +
			`"BatchId":5,` +
			`"Stuff-2":"Middle",` +
			`"RecordId":6,` +
			`"ImageMask":"def",` +
			`"Stuff-3":"End"` +
			`}]}`)
	got, err := ParseJSON(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tabulatorID0, batchID0, recordID0, session0Err := got.GetSessionIDs(0)
	imageMask0, image0Err := got.GetSessionImageMask(0)
	tabulatorID1, batchID1, recordID1, session1Err := got.GetSessionIDs(1)
	imageMask1, image1Err := got.GetSessionImageMask(1)
	tabulatorID2, batchID2, recordID2, session2Err := got.GetSessionIDs(2)
	imageMask2, image2Err := got.GetSessionImageMask(2)

	switch {
	case got.NumSessions() != 2:
		t.Errorf("incorrect number of sessions")
	case got.GetVersionString() != "4.5.6":
		t.Errorf("incorrect version string")
	case session0Err != nil:
		t.Errorf("incorrect session 0 fetch")
	case tabulatorID0 != 1:
		t.Errorf("incorrect session 0 -- tabulator ID")
	case batchID0 != 2:
		t.Errorf("incorrect session 0 -- batch ID")
	case recordID0 != 3:
		t.Errorf("incorrect session 0 -- record ID")
	case image0Err != nil:
		t.Errorf("incorrect image-mask 0 fetch")
	case imageMask0 != "ABC":
		t.Errorf("incorrect image-mask 0")
	case session1Err != nil:
		t.Errorf("incorrect session 1 fetch")
	case tabulatorID1 != 4:
		t.Errorf("incorrect session 1 -- tabulator ID")
	case batchID1 != 5:
		t.Errorf("incorrect session 1 -- batch ID")
	case recordID1 != 6:
		t.Errorf("incorrect session 1 -- record ID")
	case image1Err != nil:
		t.Errorf("incorrect image-mask 1 fetch")
	case imageMask1 != "def":
		t.Errorf("incorrect image-mask 1")
	case session2Err == nil:
		t.Errorf("incorrect session 2 fetch")
	case tabulatorID2 != -1:
		t.Errorf("incorrect no-session return -- tabulator ID")
	case batchID2 != -1:
		t.Errorf("incorrect no-session return -- batch ID")
	case recordID2 != -1:
		t.Errorf("incorrect no-session return -- record ID")
	case image2Err == nil:
		t.Errorf("incorrect image-mask 2 fetch")
	case imageMask2 != "":
		t.Errorf("incorrect no-session return -- image mask")
	}
}

func Test_JSONBlob_SetValues(t *testing.T) {
	// This is too much under-test in a single testcase but whole lot of
	// boiler-plat to break-up
	input := []byte(
		`{"Version":"4.5.6","Sessions":[{` +
			`"TabulatorId":1,` +
			`"BatchId":2,` +
			`"RecordId":3,` +
			`"ImageMask":"ABC"` +
			`},{` +
			`"Stuff-1":"Begin",` +
			`"TabulatorId":4,` +
			`"BatchId":5,` +
			`"Stuff-2":"Middle",` +
			`"RecordId":6,` +
			`"ImageMask":"def",` +
			`"Stuff-3":"End"` +
			`}]}`)

	// Golang's json module guarantees ordering of key so safe to have this be a constant.
	// https://pkg.go.dev/encoding/json#Marshal
	want := []byte(
		`{"Sessions":[{` +
			`"BatchId":2,` +
			`"ImageMask":"XXXX",` +
			`"RecordId":"0x1111",` +
			`"TabulatorId":1` +
			`},{` +
			`"BatchId":5,` +
			`"ImageMask":"YYYY",` +
			`"RecordId":"0x2222",` +
			`"Stuff-1":"Begin",` +
			`"Stuff-2":"Middle",` +
			`"Stuff-3":"End",` +
			`"TabulatorId":4` +
			`}],` +
			`"Version":"4.5.6"}`)

	object, err := ParseJSON(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := object.SetSessionRecordID(0, "0x1111"); err != nil {
		t.Fatalf("error setting 0's record ID")
	}
	if err := object.SetSessionImageMask(0, "XXXX"); err != nil {
		t.Fatalf("error setting 0's image mask")
	}
	if err := object.SetSessionRecordID(1, "0x2222"); err != nil {
		t.Fatalf("error setting 1's record ID")
	}
	if err := object.SetSessionImageMask(1, "YYYY"); err != nil {
		t.Fatalf("error setting 1's image mask")
	}

	got, err := GetJSONBytes(object, []string{"0x1111", "0x2222"}) // No re-ordering required.
	if err != nil {
		t.Fatalf("unexpected marshal error: %v", err)
	}

	if !bytes.Equal(got, want) {
		fmt.Println(string(got))
		fmt.Println(string(want))
		t.Errorf("marshalled bytes are wrong: %s", string(got))
	}
}

func Test_JSONBlob_ParseError(t *testing.T) {
	testcases := []struct {
		desc    string
		input   string
		wantErr string
	}{
		{
			desc:    "invalid JSON",
			input:   `}{`,
			wantErr: "root-unmarshal",
		},
		{
			desc:    "invalid session list",
			input:   `{"Sessions":[}`,
			wantErr: "root-unmarshal",
		},
		{
			desc:    "invalid session list contents",
			input:   `{"Sessions":[aaaaa]}`,
			wantErr: "root-unmarshal",
		},
		{
			desc:    "missing Session key",
			input:   `{}`,
			wantErr: "'Sessions' key",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.desc, func(t *testing.T) {
			got, err := ParseJSON([]byte(tc.input))
			if err == nil {
				t.Fatalf("unexpected success")
			}
			if !strings.Contains(err.Error(), tc.wantErr) {
				t.Errorf("incorrect error: %v", err)
			}
			if got != nil {
				t.Errorf("non-nil return on error")
			}
		})
	}

}

func Test_ReorderingSessions(t *testing.T) {
	testcases := []struct {
		// json.RawMessage is just a type re-cast from []byte
		// https://pkg.go.dev/encoding/json#RawMessage
		desc       string
		inputList  []string
		inputIDs   []string
		reorderIDs []string
		want       []string
	}{
		{
			desc:       "already in-order",
			inputList:  []string{"aaa", "bbb", "ccc"},
			inputIDs:   []string{"1", "2", "3"},
			reorderIDs: []string{"1", "2", "3"},
			want:       []string{"aaa", "bbb", "ccc"},
		},
		{
			desc:       "reverse order",
			inputList:  []string{"aaa", "bbb", "ccc"},
			inputIDs:   []string{"3", "2", "1"},
			reorderIDs: []string{"1", "2", "3"},
			want:       []string{"ccc", "bbb", "aaa"},
		},
		{
			desc:       "other order",
			inputList:  []string{"aaa", "bbb", "ccc"},
			inputIDs:   []string{"3", "1", "2"},
			reorderIDs: []string{"1", "2", "3"},
			want:       []string{"bbb", "ccc", "aaa"},
		},
		{
			desc:       "not based on contents",
			inputList:  []string{"8iI8", "5KR8", "edZW", "5RGa", "jmja"},
			inputIDs:   []string{"yH1T", "fh1S", "ULBO", "MwIJ", "8xBd"},
			reorderIDs: []string{"ULBO", "8xBd", "yH1T", "MwIJ", "fh1S"},
			want:       []string{"edZW", "jmja", "8iI8", "5RGa", "5KR8"},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.desc, func(t *testing.T) {
			realWant := make([]json.RawMessage, len(tc.inputList), len(tc.inputList))
			for i, value := range tc.want {
				realWant[i] = json.RawMessage([]byte(value))
			}
			realInput := make([]json.RawMessage, len(tc.inputList), len(tc.inputList))
			for i, value := range tc.inputList {
				realInput[i] = json.RawMessage([]byte(value))
			}

			got := reorderSessionList(realInput, tc.inputIDs, tc.reorderIDs)
			if !reflect.DeepEqual(got, realWant) {
				t.Errorf("incorrect output")
			}
		})
	}
}

func Test_JSONBlob_VersionString(t *testing.T) {
	blob, err := ParseJSON([]byte(`{"Version":"1.2.3","Sessions":[]}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := blob.GetVersionString()
	if got != "1.2.3" {
		t.Errorf("incorrect output. got: %q, want: '1.2.3'", got)
	}
}

func Test_JSONBlob_NoVersionString(t *testing.T) {
	blob, err := ParseJSON([]byte(`{"Sessions":[]}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := blob.GetVersionString()
	if got != "NOT-AVAILABLE" {
		t.Errorf("incorrect output. got: %q, want: 'NOT-AVAILABLE'", got)
	}
}

func Test_JSONBlob_GetSessionError(t *testing.T) {
	testcases := []struct {
		desc    string
		input   []byte
		wantErr string
	}{
		{
			desc: "Tabulator ID is JSON-string",
			input: []byte(`{"Version":"4.5.6","Sessions":[{` +
				`"TabulatorId":"1",` +
				`"BatchId":2,` +
				`"RecordId":3,` +
				`"ImageMask":"ABC"` +
				`}]}`),
			wantErr: "type int",
		},
		{
			desc: "Batch ID is JSON-string",
			input: []byte(`{"Version":"4.5.6","Sessions":[{` +
				`"TabulatorId":1,` +
				`"BatchId":"2",` +
				`"RecordId":3,` +
				`"ImageMask":"ABC"` +
				`}]}`),
			wantErr: "type int",
		},
		{
			desc: "Record ID is JSON-string",
			input: []byte(`{"Version":"4.5.6","Sessions":[{` +
				`"TabulatorId":1,` +
				`"BatchId":2,` +
				`"RecordId":"3",` +
				`"ImageMask":"ABC"` +
				`}]}`),
			wantErr: "type int",
		},
		{
			desc: "Tabulator ID is missing",
			input: []byte(
				`{"Version":"4.5.6","Sessions":[{` +
					`"BatchId":2,` +
					`"RecordId":3,` +
					`"ImageMask":"ABC"` +
					`}]}`),
			wantErr: "unexpected end",
		},
		{
			desc: "Batch ID is missing",
			input: []byte(
				`{"Version":"4.5.6","Sessions":[{` +
					`"TabulatorId":1,` +
					`"RecordId":3,` +
					`"ImageMask":"ABC"` +
					`}]}`),
			wantErr: "unexpected end",
		},
		{
			desc: "Record ID is missing",
			input: []byte(
				`{"Version":"4.5.6","Sessions":[{` +
					`"TabulatorId":1,` +
					`"BatchId":2,` +
					`"ImageMask":"ABC"` +
					`}]}`),
			wantErr: "unexpected end",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.desc, func(t *testing.T) {

			blob, err := ParseJSON(tc.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			_, _, _, err = blob.GetSessionIDs(0)
			if !strings.Contains(err.Error(), tc.wantErr) {
				t.Errorf("incorrect error. got: %v, want-with: %q'", err, tc.wantErr)
			}
		})
	}
}
