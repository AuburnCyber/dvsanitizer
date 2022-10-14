package filefmt

import (
	"fmt"
	"reflect"
	"testing"
)

func Test_OddballCSVFormat(t *testing.T) {
	testcases := []struct {
		inputLines [][]byte
		want       [][]string
	}{
		{
			inputLines: [][]byte{
				[]byte(`"a","b","c"`),
			},
			want: [][]string{
				[]string{"a", "b", "c"},
			},
		},
		{
			inputLines: [][]byte{
				[]byte(`"a",="b",="c"`),
			},
			want: [][]string{
				[]string{"a", "b", "c"},
			},
		},
		{
			inputLines: [][]byte{
				[]byte(`"a",="b",="c"`),
				[]byte(`="d",="e",="f"`),
			},
			want: [][]string{
				[]string{"a", "b", "c"},
				[]string{"d", "e", "f"},
			},
		},
	}

	for i, tc := range testcases {
		t.Run(fmt.Sprintf("TC: #%d", i), func(t *testing.T) {
			input := []byte{}
			for _, line := range tc.inputLines {
				input = append(input, []byte(line)...)
				input = append(input, '\r')
				input = append(input, '\n')
			}
			got, gotOddball, err := ReadCSVBytes(input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("incorrect output: %v", got)
			}
			if !gotOddball {
				t.Errorf("unexpected standard CSV format indicated")
			}

		})
	}
}

func Test_StandardCSVFormat(t *testing.T) {
	testcases := []struct {
		inputLines [][]byte
		want       [][]string
	}{
		{
			inputLines: [][]byte{
				[]byte(`a,b,c`),
			},
			want: [][]string{
				[]string{"a", "b", "c"},
			},
		},
		{
			inputLines: [][]byte{
				[]byte(`a,b,c`),
				[]byte(`d,e,f`),
			},
			want: [][]string{
				[]string{"a", "b", "c"},
				[]string{"d", "e", "f"},
			},
		},
	}

	for i, tc := range testcases {
		t.Run(fmt.Sprintf("TC: #%d", i), func(t *testing.T) {
			input := []byte{}
			for _, line := range tc.inputLines {
				input = append(input, []byte(line)...)
				input = append(input, '\r')
				input = append(input, '\n')
			}
			got, gotOddball, err := ReadCSVBytes(input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("incorrect output: %v", got)
			}
			if gotOddball {
				t.Errorf("unexpected oddball CSV format indicated")
			}

		})
	}
}
