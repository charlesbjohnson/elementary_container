package filedigest

import (
	"encoding/hex"
	"hash"
	"io"
	"os"
)

func Digest(file *os.File, strategy hash.Hash) (string, error) {
	if _, err := io.Copy(strategy, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(strategy.Sum(nil)), nil
}
