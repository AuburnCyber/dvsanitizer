package filefmt

import (
	"strconv"
	"strings"
	"testing"
)

func Test_LinesToCSVBytes_OddballOK(t *testing.T) {
	testcases := []struct {
		line []string
		want string
	}{
		{
			line: []string{"", "", "", "", "", "", ""},
			want: `="",="",="",="",="","",""` + "\r\n",
		},
		{
			line: []string{"0", "1", "2", "3", "4", "5", "6"},
			want: `="0",="1",="2",="3",="4","5","6"` + "\r\n",
		},
		{
			line: []string{"000", "01", "020", "", "", "", ""},
			want: `="000",="01",="020",="",="","",""` + "\r\n",
		},
		{
			line: []string{"a", "b", "c", "", "", "", ""},
			want: `"a","b","c",="",="","",""` + "\r\n",
		},
		{
			line: []string{"a", "b", "c", "3", "4", "5", "6"},
			want: `"a","b","c",="3",="4","5","6"` + "\r\n",
		},
		{
			line: []string{"a", "0", "c", "", "", "", ""},
			want: `"a",="0","c",="",="","",""` + "\r\n",
		},
		{
			line: []string{"ab", `c""d`, "ef", "", "", "", ""},
			want: `"ab","c""d","ef",="",="","",""` + "\r\n",
		},
		{
			line: []string{"0", "1", "2", "3", "4-4-4", "", ""},
			want: `="0",="1",="2",="3",="4-4-4","",""` + "\r\n",
		},
	}
	marshalledHeaderLines := `"","","","","","",""` + "\r\n" +
		`"","","","","","",""` + "\r\n" +
		`"","","","","","",""` + "\r\n" +
		`"","","","","","",""` + "\r\n"

	for i, tc := range testcases {
		t.Run(strconv.Itoa(i)+" ## "+strings.Join(tc.line, " --- "), func(t *testing.T) {
			got, err := LinesToCSVBytes([][]string{
				[]string{"", "", "", "", "", "", ""},
				[]string{"", "", "", "", "", "", ""},
				[]string{"", "", "", "", "", "", ""},
				[]string{"", "", "", "", "", "", ""},
				tc.line}, true)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// Ensure header-lines are correct
			if !strings.HasPrefix(string(got), marshalledHeaderLines) {
				t.Fatalf("marshalled fake-header lines are incorrect. got: %s\n\nwant: %s", string(got), marshalledHeaderLines)
			}
			got = got[len(marshalledHeaderLines):]

			// Check the data-line
			if string(got) != tc.want {
				t.Errorf("incorrect output. got: %q, want: %q", string(got), tc.want)
			}
		})
	}
}
