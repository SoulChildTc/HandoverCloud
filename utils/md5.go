package utils

import (
	"crypto/md5"
	"encoding/hex"
)

func PasswdMd5Digest(txt string) string {
	if txt == "" {
		return ""
	}

	t := []byte(txt)

	h := md5.New()
	h.Write(t[len(t)-1:])
	h.Write(t)
	h.Write(t[0:1])

	return hex.EncodeToString(h.Sum(nil)[:])
}

func Md5Digest(txt string) string {
	h := md5.New()
	h.Write([]byte(txt))
	return hex.EncodeToString(h.Sum(nil)[:])
}
