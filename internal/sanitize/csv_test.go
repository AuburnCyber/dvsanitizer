package sanitize

import (
	"fmt"
	"reflect"
	"testing"
)

func Test_CSVDataLineSorting(t *testing.T) {
	testcases := []struct {
		topDataLine    []string
		bottomDataLine []string
		wantFlipped    bool
	}{
		{[]string{"0", "1", "1", "1"}, []string{"0", "1", "1", "2"}, false},    // In-order
		{[]string{"0", "1", "1", "2"}, []string{"0", "1", "1", "1"}, true},     // Out-of-order
		{[]string{"0", "1", "2", "2"}, []string{"0", "1", "2", "1"}, true},     // Batch difference
		{[]string{"0", "2", "1", "2"}, []string{"0", "2", "1", "1"}, true},     // Tabulator difference
		{[]string{"0", "1", "1", "10000"}, []string{"0", "1", "1", "9"}, true}, // Padding needed
	}

	for i, tc := range testcases {
		t.Run(fmt.Sprintf("TC #%d", i), func(t *testing.T) {
			input := [][]string{tc.topDataLine, tc.bottomDataLine}

			got := resortCSVDataLines(input)

			var wantTop []string
			var wantBottom []string
			if tc.wantFlipped {
				wantTop = tc.bottomDataLine
				wantBottom = tc.topDataLine
			} else {
				wantTop = tc.topDataLine
				wantBottom = tc.bottomDataLine
			}

			if !reflect.DeepEqual(got[0], wantTop) {
				t.Errorf("incorrect top-line")
			}
			if !reflect.DeepEqual(got[1], wantBottom) {
				t.Errorf("incorrect bottom-line")
			}
			if got[0][0] != "" {
				t.Errorf("non-empty CVR ID")
			}
			if got[1][0] != "" {
				t.Errorf("non-empty CVR ID")
			}
		})
	}
}
