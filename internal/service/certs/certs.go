package certs

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/asn1"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"
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
	dir string
}

func New() (*Certs, error) {
	if _, err := os.Stat(certsDir); os.IsNotExist(err) {
		if err = os.Mkdir(certsDir, certsDirPerm); err != nil {
			return nil, err
		}
	}
	return &Certs{dir: certsDir}, nil
}

func (c *Certs) RsaKeys() (public, private []byte, err error) {
	_, pubErr := os.Stat(c.filePath(publicFile))
	_, prvErr := os.Stat(c.filePath(privateFile))

	if os.IsNotExist(pubErr) || os.IsNotExist(prvErr) {
		if err = c.createRsaCerts(); err != nil {
			return nil, nil, fmt.Errorf("fail create rsa certs %w", err)
		}
	}

	pubKey, err := os.ReadFile(c.filePath(publicFile))
	if err != nil {
		return nil, nil, fmt.Errorf("fail read public rsa cert %w", err)
	}

	prvKey, err := os.ReadFile(c.filePath(privateFile))
	if err != nil {
		return nil, nil, fmt.Errorf("fail read private rsa cert %w", err)
	}

	return pubKey, prvKey, nil
}

func (c *Certs) RemoveDir() error {
	return os.RemoveAll(c.dir)
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
