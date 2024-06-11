package crypt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	testCases := []struct {
		name    string
		secret  string
		isPanic bool
	}{
		{
			name:    "Secret invalid",
			secret:  "1234567890",
			isPanic: true,
		},
		{
			name:    "Secret AES128",
			secret:  "1234567890123456",
			isPanic: false,
		},
		{
			name:    "Secret AES192",
			secret:  "123456789012345678901234",
			isPanic: false,
		},
		{
			name:    "Secret AES256",
			secret:  "12345678901234567890123456789012",
			isPanic: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.isPanic {
				assert.Panics(t, func() {
					New(tc.secret)
				})
			} else {
				assert.NotPanics(t, func() {
					New(tc.secret)
				})
			}
		})
	}
}

func TestService_Encrypt_Decrypt(t *testing.T) {
	testCases := []struct {
		name  string
		input string
	}{
		{
			name:  "Empty text",
			input: "",
		},
		{
			name:  "Alpha text",
			input: "some text",
		},
		{
			name:  "Number text",
			input: "123 456 789",
		},
		{
			name:  "Symbol text",
			input: "~!@#$%*",
		},
	}

	crypt := New("1234567890123456")

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			hash, err := crypt.Encrypt(tc.input)
			assert.NoError(t, err)

			text, err := crypt.Decrypt(hash)
			assert.NoError(t, err)

			assert.Equal(t, tc.input, text)
		})
	}
}

func TestService_HashPassword(t *testing.T) {
	testCases := []struct {
		name string
		pass string
	}{
		{
			name: "Empty password",
			pass: "",
		},
		{
			name: "Text password",
			pass: "some text",
		},
		{
			name: "Alpha password",
			pass: "qwerty",
		},
		{
			name: "Number password",
			pass: "123456",
		},
		{
			name: "Alpha numeric unicode password",
			pass: "qwerty123456~@$%&*",
		},
	}

	crypt := New("1234567890123456")

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			hash, err := crypt.HashPassword(tc.pass)
			assert.NoError(t, err)
			assert.True(t, crypt.CompareHashPassword(tc.pass, hash))
		})
	}
}
