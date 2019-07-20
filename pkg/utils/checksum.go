package utils

import (
	"crypto/md5"
	"encoding/hex"
	"io"
)

func Md5Checksum(r io.Reader) (string, error) {
	hash := md5.New()

	if _, err := io.Copy(hash, r); err != nil {
		return "", err
	}

	hashInBytes := hash.Sum(nil)[:16]
	return hex.EncodeToString(hashInBytes), nil
}