package certs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKeys(t *testing.T) {
	certs, err := New()
	assert.NoError(t, err)

	public, private, err := certs.Keys()
	assert.NoError(t, err)
	assert.NotNil(t, public)
	assert.NotNil(t, private)

	public2, err := certs.PublicKey()
	assert.NoError(t, err)
	assert.Equal(t, public, public2)

	private2, err := certs.PrivateKey()
	assert.NoError(t, err)
	assert.Equal(t, private, private2)

	err = certs.RemoveDir()
	assert.NoError(t, err)
}

func TestJWK(t *testing.T) {
	certs, err := New()
	assert.NoError(t, err)

	public, err := certs.PublicKey()
	assert.NoError(t, err)
	assert.NotNil(t, public)

	jwk, err := certs.PublicJWK()
	assert.NoError(t, err)

	public2, err := certs.PublicByJWK(jwk)
	assert.NoError(t, err)

	assert.Equal(t, public, public2)

	err = certs.RemoveDir()
	assert.NoError(t, err)
}
