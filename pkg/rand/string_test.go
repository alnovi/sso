package rand

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const StringLength = 20

func TestBase64(t *testing.T) {
	str1 := Base64(StringLength)
	str2 := Base64(StringLength)

	assert.NotEqual(t, str1, str2)
	assert.Len(t, str1, StringLength)
	assert.Len(t, str2, StringLength)
}

func TestBase62(t *testing.T) {
	str1 := Base62(StringLength)
	str2 := Base62(StringLength)

	assert.NotEqual(t, str1, str2)
	assert.Len(t, str1, StringLength)
	assert.Len(t, str2, StringLength)
}

func TestHex(t *testing.T) {
	str1 := Hex(StringLength)
	str2 := Hex(StringLength)

	assert.NotEqual(t, str1, str2)
	assert.Len(t, str1, StringLength)
	assert.Len(t, str2, StringLength)
}

func TestDec(t *testing.T) {
	str1 := Dec(StringLength)
	str2 := Dec(StringLength)

	assert.NotEqual(t, str1, str2)
	assert.Len(t, str1, StringLength)
	assert.Len(t, str2, StringLength)
}

func TestString(t *testing.T) {
	chars := Charset("абвгдеёжзийклмнопрстуфхцчшщъыьэюяАБВГДЕЁЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯ")

	str1 := String(StringLength, chars)
	str2 := String(StringLength, chars)

	assert.NotEqual(t, str1, str2)
	assert.Len(t, str1, StringLength)
	assert.Len(t, str2, StringLength)
}
