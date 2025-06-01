package certs

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/asn1"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"path/filepath"

	"github.com/golang-jwt/jwt/v5"
)

const (
	certsDir     = "./certs"
	certsDirPerm = 0700
	certPerm     = 0600
	rsaBits      = 2048
	privateFile  = "private.pem"
	publicFile   = "public.pem"
)

type Certs struct {
	dir     string
	public  *rsa.PublicKey
	private *rsa.PrivateKey
}

func New() (*Certs, error) {
	if _, err := os.Stat(certsDir); os.IsNotExist(err) {
		if err = os.Mkdir(certsDir, certsDirPerm); err != nil {
			return nil, err
		}
	}

	certs := &Certs{dir: certsDir}

	return certs, certs.initCerts()
}

func (c *Certs) Keys() (*rsa.PublicKey, *rsa.PrivateKey, error) {
	pubKey, err := c.PublicKey()
	if err != nil {
		return nil, nil, err
	}

	prvKey, err := c.PrivateKey()
	if err != nil {
		return nil, nil, err
	}

	return pubKey, prvKey, nil
}

func (c *Certs) PublicKey() (*rsa.PublicKey, error) {
	if c.public != nil {
		return c.public, nil
	}

	pubKey, err := os.ReadFile(c.filePath(publicFile))
	if err != nil {
		return nil, fmt.Errorf("fail read public rsa cert %w", err)
	}

	c.public, err = jwt.ParseRSAPublicKeyFromPEM(pubKey)

	return c.public, err
}

func (c *Certs) PrivateKey() (*rsa.PrivateKey, error) {
	if c.private != nil {
		return c.private, nil
	}

	prvKey, err := os.ReadFile(c.filePath(privateFile))
	if err != nil {
		return nil, fmt.Errorf("fail read private rsa cert %w", err)
	}

	c.private, err = jwt.ParseRSAPrivateKeyFromPEM(prvKey)

	return c.private, err
}

func (c *Certs) PublicJWK() (*JWK, error) {
	key, err := c.PublicKey()
	if err != nil {
		return nil, err
	}
	return NewJwk(key), nil
}

func (c *Certs) PublicByJWK(jwk *JWK) (*rsa.PublicKey, error) {
	if jwk == nil {
		return nil, fmt.Errorf("jwk is nil")
	}

	n, err := c.base64ToBigInt(jwk.N)
	if err != nil {
		return nil, err
	}

	e, err := c.base64ToBigInt(jwk.E)
	if err != nil {
		return nil, err
	}

	return &rsa.PublicKey{N: n, E: int(e.Int64())}, nil
}

func (c *Certs) RemoveDir() error {
	err := os.RemoveAll(c.dir)
	if err != nil {
		return err
	}

	c.public = nil
	c.private = nil
	return nil
}

func (c *Certs) initCerts() error {
	_, pubErr := os.Stat(c.filePath(publicFile))
	_, prvErr := os.Stat(c.filePath(privateFile))

	if os.IsNotExist(pubErr) || os.IsNotExist(prvErr) {
		if err := c.createRsaCerts(); err != nil {
			return fmt.Errorf("fail create rsa certs %w", err)
		}
	}

	return nil
}

func (c *Certs) filePath(name string) string {
	return filepath.Join(c.dir, name)
}

func (c *Certs) createRsaCerts() error {
	buf := bytes.NewBuffer(nil)

	private, err := rsa.GenerateKey(rand.Reader, rsaBits)
	if err != nil {
		return fmt.Errorf("fail generate private.pem: %w", err)
	}

	privateKey := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(private),
	}

	err = pem.Encode(buf, privateKey)
	if err != nil {
		return fmt.Errorf("fail encode private.pem: %w", err)
	}

	if err = os.WriteFile(c.filePath(privateFile), buf.Bytes(), certPerm); err != nil {
		return fmt.Errorf("cannot write private.pem: %w", err)
	}

	buf = bytes.NewBuffer(nil)

	public, err := asn1.Marshal(private.PublicKey)
	if err != nil {
		return fmt.Errorf("fail marshal RSA key: %w", err)
	}

	var pemKey = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: public,
	}

	err = pem.Encode(buf, pemKey)
	if err != nil {
		return fmt.Errorf("fail encode public.pem: %w", err)
	}

	if err = os.WriteFile(c.filePath(publicFile), buf.Bytes(), certPerm); err != nil {
		panic(fmt.Errorf("cannot write public.pem: %s\n", err))
	}

	return nil
}

func (c *Certs) base64ToBigInt(val string) (*big.Int, error) {
	b, err := base64.URLEncoding.DecodeString(val)
	return new(big.Int).SetBytes(b), err
}
