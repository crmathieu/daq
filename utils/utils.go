package utils

import (
	"encoding/hex"
	"fmt"
)

func Decode(hs string) string {
	decoded, err := hex.DecodeString(hs)
	if err != nil {
		fmt.Printf("hex.DecodeString - %s\n", err.Error())
	}
	return string(decoded)
}

func Encode(s string) string {
	encoded := hex.EncodeToString([]byte(s))
	return encoded
}
