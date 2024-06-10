package rand

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const length = 10

func TestBase64(t *testing.T) {
	str := Base64(length)
	assert.Equal(t, length, len(str))
}

func TestBase62(t *testing.T) {
	str := Base62(length)
	assert.Equal(t, length, len(str))
}

func TestHex(t *testing.T) {
	hex := Hex(length)
	assert.Equal(t, length, len(hex))
}

func TestDec(t *testing.T) {
	dec := Dec(length)
	assert.Equal(t, length, len(dec))
}

func TestString(t *testing.T) {
	str := String(length, "1")
	assert.Equal(t, length, len(str))
	assert.Equal(t, "1111111111", str)
}

func TestStringLengthPanic(t *testing.T) {
	assert.Panics(t, func() {
		_ = String(-1, Base64Chars)
	})
}

func TestStringCharsPanic(t *testing.T) {
	assert.Panics(t, func() {
		_ = String(length, "")
	})
}
