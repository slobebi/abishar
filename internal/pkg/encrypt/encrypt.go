package encrypt

import (
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha512"
	"fmt"
	"io"

	"github.com/gofrs/uuid"
)

var combination = []byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

// GenerateUUID generate new UUID
func GenerateUUID() (string, error) {
	generated, err := uuid.NewV4()
	if err != nil {
		return "", err
	}

	return generated.String(), nil
}

// EncodeSHA1 encode string using SHA1. Return as hex
func EncodeSHA1(str string) string {
	h := sha1.New()
	h.Write([]byte(str))
	bs := h.Sum(nil)

	return fmt.Sprintf("%x", bs)
}

// EncodeSHA512 encode string using SHA512. Return as hex
func EncodeSHA512(str string) string {
	h := sha512.New()
	h.Write([]byte(str))
	bs := h.Sum(nil)

	return fmt.Sprintf("%x", bs)
}

// RandomizeNumber encode string. return combination of 4 character
func RandomizeNumber(length int) string {

	b := make([]byte, length)
	n, err := io.ReadAtLeast(rand.Reader, b, length)
	if n != length || err != nil {
		return ""
	}
	for i := 0; i < len(b); i++ {
		b[i] = combination[int(b[i])%len(combination)]
	}
	return string(b)
}
