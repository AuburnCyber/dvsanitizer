package sanitize

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"

	"github.com/AuburnCyber/dvsanitizer/internal/safety"
)

const (
	ONE_MILLION = 1000000
)

var aesCipher cipher.Block

func Initialize(seed []byte) {
	// Actual checks occur in flags.go's ParseCmd() of the main app.
	if len(seed) < 16 {
		// This is a last-ditch to avoid stupidity.
		panic("CODE PATH ALLOWED CALL TO sanitize.Initialize() WITH INSECURE SEED")
	}
	keyArr := sha256.Sum256(seed)

	var err error
	aesCipher, err = aes.NewCipher(keyArr[:])
	if err != nil {
		// Don't report anything else b/c don't want to leak a seed that may be reused.
		safety.ReportError("Unable to create AES cipher object.", err)
	}
}

func createCleanRecordID(tabID uint32, batchID uint32, dirtyRecordID uint32) string {
	pt := make([]byte, aes.BlockSize, aes.BlockSize)
	binary.BigEndian.PutUint32(pt, uint32(tabID))
	binary.BigEndian.PutUint32(pt[4:], uint32(batchID))
	binary.BigEndian.PutUint32(pt[8:], uint32(dirtyRecordID))
	binary.BigEndian.PutUint32(pt[12:], uint32(0)) // Ensures the last 4-bytes are null

	ct := make([]byte, aes.BlockSize, aes.BlockSize)

	// Use ECB mode which is safe in this case b/c we explicitly want identical
	// plaintexts (the record IDs) to have identical ciphertext (the encrypted
	// record IDs).
	aesCipher.Encrypt(ct, pt)

	// The sanitized record ID is the hex-encoded first 8 bytes of the
	// ciphertext prefixed with '0x' to unequivocally indicate that it is a hex
	// value and not an integer value.
	return "0x" + hex.EncodeToString(ct[:8])
}
