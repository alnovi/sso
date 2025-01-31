package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

const (
	AES128 = 16
	AES192 = 24
	AES256 = 32

	passwordCost = 14
)

var (
	ErrInvalidSecret = errors.New("crypt secret invalid size")
)

type Crypt struct {
	secret []byte
}

func New(secret string) (*Crypt, error) {
	if len(secret) != AES128 && len(secret) != AES192 && len(secret) != AES256 {
		return nil, ErrInvalidSecret
	}
	return &Crypt{secret: []byte(secret)}, nil
}

func (s *Crypt) Encrypt(input string) (string, error) {
	block, err := aes.NewCipher(s.secret)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = rand.Read(nonce); err != nil {
		return "", err
	}

	text := gcm.Seal(nonce, nonce, []byte(input), s.secret)

	return base64.RawURLEncoding.EncodeToString(text), nil
}

func (s *Crypt) Decrypt(input string) (string, error) {
	cipherText, err := base64.RawURLEncoding.DecodeString(input)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(s.secret)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	nonce, cipherText := cipherText[:nonceSize], cipherText[nonceSize:]

	text, err := gcm.Open(nil, nonce, cipherText, s.secret)
	if err != nil {
		return "", err
	}

	return string(text), nil
}

func (s *Crypt) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), passwordCost)
	return string(hash), err
}

func (s *Crypt) CompareHashPassword(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
