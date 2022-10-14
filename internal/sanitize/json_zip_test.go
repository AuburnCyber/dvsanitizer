package sanitize

import (
	"fmt"
	"testing"
)

func Test_CleanImageMask_OK(t *testing.T) {
	testcases := []struct {
		inputMask     string
		expectDirtyID int
		cleanID       string
		want          string
	}{
		{
			inputMask:     `C:\NAS\abc\Results\Tabulator00001\Batch001\Images\00001_00001_000001*.*`,
			expectDirtyID: 1,
			cleanID:       "0x1111111111111111",
			want:          `C:\NAS\abc\Results\Tabulator00001\Batch001\Images\00001_00001_0x1111111111111111*.*`,
		},
		{
			inputMask:     `C:\NAS\abc\Results\Tabulator00001\Batch001\Images\00001_00001_000001*.*`,
			expectDirtyID: 1,
			cleanID:       "0x2222222222222222",
			want:          `C:\NAS\abc\Results\Tabulator00001\Batch001\Images\00001_00001_0x2222222222222222*.*`,
		},
		{
			inputMask:     `C:\NAS\abc\Results\Tabulator00001\Batch001\Images\00001_00001_099999*.*`,
			expectDirtyID: 99999,
			cleanID:       "0x3333333333333333",
			want:          `C:\NAS\abc\Results\Tabulator00001\Batch001\Images\00001_00001_0x3333333333333333*.*`,
		},
		{
			inputMask:     `C:\NAS\abc\Results\Tabulator00001\Batch001\Images\00001_00001_000001*.*`,
			expectDirtyID: 1,
			cleanID:       "0x4444444444444444",
			want:          `C:\NAS\abc\Results\Tabulator00001\Batch001\Images\00001_00001_0x4444444444444444*.*`,
		},
		{
			inputMask:     `D:\NAS\abc\Results\Tabulator00556\Batch000\Images\00556_00000_2339149_*.*`,
			expectDirtyID: 2339149,
			cleanID:       "0x5555555555555555",
			want:          `D:\NAS\abc\Results\Tabulator00556\Batch000\Images\00556_00000_0x5555555555555555_*.*`,
		},
	}

	for _, tc := range testcases {
		t.Run(fmt.Sprintf("%d --> %s", tc.expectDirtyID, tc.cleanID), func(t *testing.T) {
			got, err := createCleanImageMask(tc.inputMask, tc.expectDirtyID, tc.cleanID)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got != tc.want {
				t.Errorf("incorrect output.\ngot: %q\nwant: %q", got, tc.want)
			}
		})
	}
}
