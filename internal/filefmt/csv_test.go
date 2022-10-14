package filefmt

import (
	"reflect"
	"testing"
)

func Test_ConvertCSVLineToFields(t *testing.T) {
	testcases := []struct {
		line string
		want []string
	}{
		{
			line: `"","",""`,
			want: []string{"", "", ""},
		},
		{
			line: `"1","2","3"`,
			want: []string{"1", "2", "3"},
		},
		{
			line: `="123",="456",="789"`,
			want: []string{"123", "456", "789"},
		},
		{
			line: `"0",="0",="0000"`,
			want: []string{"0", "0", "0000"},
		},
		{
			line: `"ab","cd","ef"`,
			want: []string{"ab", "cd", "ef"},
		},
		{
			line: `"ab","c""d","ef"`,
			want: []string{"ab", `c""d`, "ef"},
		},
	}

	for _, tc := range testcases {
		got, err := ConvertCSVLineToFields(tc.line)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			continue
		}

		if !reflect.DeepEqual(got, tc.want) {
			t.Errorf("incorrect output. got: %v, want: %v", got, tc.want)
		}
	}
}
