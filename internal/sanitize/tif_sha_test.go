package sanitize

import (
	"testing"
)

func Test_IsTIFSHAFilename(t *testing.T) {
	testcases := []struct {
		input string
		want  bool
	}{
		{"00001_00002_00003.tif", true},
		{"00001_00002_00003.sha", true},
		{"abc.tif", false},
		{"00001_00002_0003.tif", false}, // Record ID missing digit
		{"00001_0002_00003.tif", false}, // Batch ID missing digit
		{"0001_00002_00003.tif", false}, // Tabulator ID missing digit
		{"no-dots", false},
		{"should.not.have.multiple.dots.tif", false},
		{"wrong_extension_on_file.fit", false},
		{"less_unders.tif", false},
		{"far_to_many_underscores_to_be_right.tif", false},
		{"this/is/a/linux/path/file.tif", false},
		{`this\is\a\windows\path\file.tif`, false},
		{"00lll_00000_2222222_3.tif", true},
		{"Thumbs.db", false},
		{"1_22_33_4_slog.txt", false},
		{"1_22_33_44_DETAIL.DVD", false},
		{"scan1111.pdf", false},
		{"Audit_Totals.xlsx", false},
		{"Audit_Totals.xls", false},
		{"TabulatorData.xml", false},
		{"11111111-2222-3333-4444-555555555555.spx", false},
		{"11111111-2222-3333-4444-555555555555.wav", false},
	}

	for _, tc := range testcases {
		t.Run(tc.input, func(t *testing.T) {

		})
	}
}

func Test_buildCleanTIFSHAPath(t *testing.T) {
	testcases := []struct {
		dirtyPath    string
		dirtyBaseDir string
		cleanBaseDir string
		want         string
	}{
		{
			dirtyPath: "1111_2222_33333.tif",
			want:      "1111_2222_0xe9a341ef3f306166.tif",
		},
		{
			dirtyPath: "1111_2222_33333.sha",
			want:      "1111_2222_0xe9a341ef3f306166.sha",
		},
		{
			dirtyPath: "1111_2222_33333_1.sha",
			want:      "1111_2222_0xe9a341ef3f306166_1.sha",
		},
		{
			dirtyPath:    "a/1111_2222_33333.tif",
			dirtyBaseDir: "a/",
			cleanBaseDir: "b/",
			want:         "b/1111_2222_0xe9a341ef3f306166.tif",
		},
		{
			dirtyPath:    "a/b/c/1111_2222_33333.tif",
			dirtyBaseDir: "a/",
			cleanBaseDir: "d/",
			want:         "d/b/c/1111_2222_0xe9a341ef3f306166.tif",
		},
		{
			dirtyPath:    "a/b/c d/1111_2222_33333.tif",
			dirtyBaseDir: "a/",
			cleanBaseDir: "d/",
			want:         "d/b/c d/1111_2222_0xe9a341ef3f306166.tif",
		},
	}
	Initialize(testingKey)

	for _, tc := range testcases {
		t.Run(tc.dirtyPath, func(t *testing.T) {
			got := buildCleanTIFSHAPath(tc.dirtyPath, tc.dirtyBaseDir, tc.cleanBaseDir)
			if got != tc.want {
				t.Errorf("incorrect output. got: %q, want: %q", got, tc.want)
			}
		})
	}
}
