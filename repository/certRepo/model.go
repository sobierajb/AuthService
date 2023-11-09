package certRepo

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
)

type KeyType int

const (
	NoValue KeyType = iota
	RSA256
	RSA384
	RSA512
)

type Certificate struct {
	Id   string
	Type KeyType
	Key  []byte
}

type PublicCertResponse struct {
	Alg string // cryptographic alghorythm used with key
	Kty string // familly of cryptographic alghorytm
	N   string //modulus for RSA public
	E   string //exponent for RSA public
	X5t string //The thumbprint
	X5c string //x.509 chain
}

func (c *Certificate) GenerateRsa(ky KeyType) error {
	var key *rsa.PrivateKey
	var err error
	switch ky {
	case RSA256:
		key, err = rsa.GenerateKey(rand.Reader, 2048)
		c.Type = RSA256
	case RSA384:
		key, err = rsa.GenerateKey(rand.Reader, 3072)
		c.Type = RSA384
	case RSA512:
		key, err = rsa.GenerateKey(rand.Reader, 4096)
		c.Type = RSA512
	default:
		return ErrWrongKeyType
	}

	if err != nil {
		return err
	}
	// Store key as String of bytes.
	marshal, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return err
	}
	c.Key = marshal
	return nil
}

func (c *Certificate) GetRsaPrivateKey() (*rsa.PrivateKey, error) {
	if c.Type != RSA256 {
		return nil, ErrWrongKeyType
	}
	key, err := x509.ParsePKCS8PrivateKey(c.Key)
	if err != nil {
		return nil, err
	}
	return key.(*rsa.PrivateKey), nil
}

func (c *Certificate) GetRsaPuiblicKey() (*rsa.PublicKey, error) {
	privateKey, err := c.GetRsaPrivateKey()
	if err != nil {
		return nil, err
	}
	return &privateKey.PublicKey, nil
}
