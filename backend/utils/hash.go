package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"mime/multipart"
)

func SHA256File(file multipart.File) (string, error) {
	h := sha256.New()
	if _, err := io.Copy(h, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
