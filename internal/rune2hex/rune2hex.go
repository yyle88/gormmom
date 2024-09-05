package rune2hex

import (
	"encoding/binary"
	"encoding/hex"
	"strings"
)

func S2LeHex4sUppers(s string) (ssHexes []string) {
	for _, c := range s {
		ssHexes = append(ssHexes, Uint32ToHex4Ups(c))
	}
	return ssHexes
}

func Uint32ToHex4Los(c int32) string {
	return strings.ToLower(Uint32ToHex4s(uint32(c)))
}

func Uint32ToHex4Ups(c int32) string {
	return strings.ToUpper(Uint32ToHex4s(uint32(c)))
}

func Uint32ToHex4s(x uint32) string {
	return hex.EncodeToString(Uint32ToBytes(x))[:4]
}

func Uint32ToBytes(x uint32) []byte {
	var src = make([]byte, 4)
	binary.LittleEndian.PutUint32(src, uint32(x))
	return src
}
