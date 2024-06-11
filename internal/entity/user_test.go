package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUser_MaskEmail(t *testing.T) {
	testCases := []struct {
		name  string
		email string
		exp   string
	}{
		{
			name:  "empty email",
			email: "",
			exp:   "",
		},
		{
			name:  "a@mail.ru",
			email: "a@mail.ru",
			exp:   "***@mail.ru",
		},
		{
			name:  "abcde@mail.ru",
			email: "abcde@mail.ru",
			exp:   "***@mail.ru",
		},
		{
			name:  "abcdef@mail.ru",
			email: "abcdef@mail.ru",
			exp:   "ab***ef@mail.ru",
		},
		{
			name:  "abcdefghijklmnop@mail.ru",
			email: "abcdefghijklmnop@mail.ru",
			exp:   "ab***op@mail.ru",
		},
		{
			name:  "abcdefghijklmnop@mail.ru",
			email: "abcdefghijklmnop@mail.ru",
			exp:   "ab***op@mail.ru",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.exp, (&User{Email: tc.email}).MaskEmail())
		})
	}
}
