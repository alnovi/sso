package rand

import (
	"crypto/rand"
	"math/big"
	"unicode/utf8"
)

type Charset string

const (
	Base64Chars Charset = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ+/"
	Base62Chars Charset = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	HexChars    Charset = "0123456789abcdef"
	DecChars    Charset = "0123456789"
)

func Base64(length int) string {
	return String(length, Base64Chars)
}

func Base62(length int) string {
	return String(length, Base62Chars)
}

func Hex(length int) string {
	return String(length, HexChars)
}

func Dec(length int) string {
	return String(length, DecChars)
}

func String(length int, charset Charset) string {
	if utf8.RuneCountInString(string(charset)) <= 0 {
		panic("character set must contain more than")
	}

	b := bytes(length)

	charsetLen := big.NewInt(int64(len(charset)))
	for i := 0; i < length; i++ {
		randomIndex, _ := rand.Int(rand.Reader, charsetLen)
		b[i] = charset[randomIndex.Int64()]
	}

	return string(b)
}

func bytes(length int) []byte {
	if length <= 0 {
		panic("length must be positive")
	}

	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return b
}
