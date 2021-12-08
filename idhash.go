package tools

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	mrand "math/rand"
	"regexp"
	"time"
)

const (
	uuidDash     byte = '-'
	randBytesMax      = 1024 * 1024
	AlgMd5            = 1
	AlgSha1           = 2
	AlgSha256         = 3
)

var UUIDRegexp = regexp.MustCompile(fmt.Sprintf("^[a-f0-9]{8}%c([a-f0-9]{4}%c){3}[a-f0-9]{12}$", uuidDash, uuidDash))

//generate random hash uuid.
func RandUUID() string {
	return uuid(randSum(16))
}

//generate hash uuid of the input data.
func HashUUID(bs []byte) string {
	return uuid(hashSum(AlgSha1, bs, 16))
}

func uuid(bs []byte) string {

	bs[6] = (bs[6] & 0x0F) | 0x40 // version 4
	bs[8] = (bs[8] & 0x3F) | 0x80 // variant rfc4122

	uuid := make([]byte, 36)

	hex.Encode(uuid[0:8], bs[0:4])
	uuid[8] = uuidDash

	hex.Encode(uuid[9:13], bs[4:6])
	uuid[13] = uuidDash

	hex.Encode(uuid[14:18], bs[6:8])
	uuid[18] = uuidDash

	hex.Encode(uuid[19:23], bs[8:10])
	uuid[23] = uuidDash

	hex.Encode(uuid[24:], bs[10:])

	return string(uuid)
}

func hashSum(alg int, bs []byte, bytelen int) []byte {
	if len(bs) == 0 {
		return bs
	}

	if bytelen < 1 {
		bytelen = 1
	}

	switch alg {
	case AlgMd5:
		if bytelen > 16 {
			bytelen = 16
		}
		hs := md5.Sum(bs)
		return hs[:bytelen]

	case AlgSha256:
		if bytelen > 32 {
			bytelen = 32
		}
		hs := sha256.Sum256(bs)
		return hs[:bytelen]

	case AlgSha1:
		if bytelen > 20 {
			bytelen = 20
		}
		hs := sha1.Sum(bs)
		return hs[:bytelen]
	}
	return []byte{}
}

func init() {
	mrand.Seed(time.Now().UTC().UnixNano())
}

var Reader io.Reader

func randSum(size int) []byte {

	if size < 1 {
		size = 1
	} else if size > randBytesMax {
		size = randBytesMax
	}

	bs := make([]byte, size)

	if _, err := rand.Read(bs); err != nil {
		for i := range bs {
			bs[i] = uint8(mrand.Intn(256))
		}
	}

	return bs
}

//check if UUID is in valid format.
func IsValidUUID(UUID string) bool {
	return UUIDRegexp.MatchString(UUID)
}
