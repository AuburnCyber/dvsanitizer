package sanitize

import (
	"fmt"
	"testing"
)

var testingKey = []byte{
	0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11,
	0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11,
	0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11,
	0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11,
}

func Test_CreateCleanRecordID(t *testing.T) {
	/*
		The following can be checked via the Python3 function:

		def sanitize(tab_id, batch_id, record_id):
			key = hashlib.sha256(bytes('\x11'*32, 'latin-1')).digest()
			cipher = AES.new(key, AES.MODE_ECB)
			pt = struct.pack('>I', tab_id) + struct.pack('>I', batch_id)
			pt += struct.pack('>I', record_id) + struct.pack('>I', 0)
			ct = cipher.encrypt(pt)
			return '0x' + ct[:8].hex()
	*/
	testcases := []struct {
		tabID    uint32
		batchID  uint32
		recordID uint32
		want     string
	}{
		{0, 0, 0, "0x4eba0cf44ba2422d"},
		{1, 1, 1, "0x058bae87a6233f8f"},
		{18, 18, 18, "0xd7c685a77d4ec7e6"},
		{31337, 31337, 31337, "0x2d5850707ece145b"},
		{ONE_MILLION - 1, ONE_MILLION - 1, ONE_MILLION - 1, "0xd9b4d8ae29223ad8"},
		{ONE_MILLION, ONE_MILLION, ONE_MILLION, "0xe4256bdd44e207f7"},
		{ONE_MILLION + 1, ONE_MILLION + 1, ONE_MILLION + 1, "0x5df6f7ca97b6b049"},
		{10*ONE_MILLION + 1, 10*ONE_MILLION + 1, 10*ONE_MILLION + 1, "0xabe42695781bb159"},
		{1111, 2222, 33333, "0xe9a341ef3f306166"},
	}
	Initialize(testingKey)

	for _, tc := range testcases {
		t.Run(fmt.Sprintf("%d -- %d -- %d", tc.tabID, tc.batchID, tc.recordID), func(t *testing.T) {
			got := createCleanRecordID(tc.tabID, tc.batchID, tc.recordID)
			if got != tc.want {
				t.Errorf("incorrect output. got: %s, want: %s", got, tc.want)
			}
		})
	}
}
