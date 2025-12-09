package tools

import (
	"crypto/md5"
	"crypto/sha256"
	"hash"
)

const (
	SignatureAlgorithmMD5        = "MD5"
	SignatureAlgorithmHMACSHA256 = "HMACSHA256"
)

func Sign(raw []byte, algo string) (signature []byte) {
	var h hash.Hash
	switch algo {
	case SignatureAlgorithmHMACSHA256:
		h = sha256.New()
	case SignatureAlgorithmMD5:
		h = md5.New()
	}
	if h != nil {
		h.Write(raw)
		signature = h.Sum(nil)
	}
	return
}
