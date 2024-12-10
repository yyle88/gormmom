package unicodehex

import (
	"encoding/binary"
	"encoding/hex"
	"strings"
)

func StringToHex4Uppercase(s string) (results []string) {
	for _, c := range s {
		results = append(results, Uint32ToHex4Uppercase(c))
	}
	return results
}

func Uint32ToHex4Lowercase(c int32) string {
	return strings.ToLower(Uint32ToHex4s(uint32(c)))
}

func Uint32ToHex4Uppercase(c int32) string {
	return strings.ToUpper(Uint32ToHex4s(uint32(c)))
}

func Uint32ToHex4s(x uint32) string {
	return hex.EncodeToString(Uint32ToLittleEndianBytes(x))[:4]
}

func Uint32ToLittleEndianBytes(x uint32) []byte {
	var src = make([]byte, 4)
	binary.LittleEndian.PutUint32(src, uint32(x))
	return src
}
