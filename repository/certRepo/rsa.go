package certRepo

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/binary"

	"github.com/golang-jwt/jwt/v4"
)

type rsaCert struct {
	data CertificateData
}

func NewRsaCert() Certificate {
	return &rsaCert{}
}

func (c *rsaCert) GetType() KeyType {
	return c.data.Type
}

func (c *rsaCert) GetRawData() *CertificateData {
	return &c.data
}

func (c *rsaCert) SetRawData(data *CertificateData) error {
	switch data.Type {
	case RSA256:
	case RSA384:
	case RSA512:
	default:
		return ErrWrongKeyType
	}
	c.data = *data
	return nil
}

func (c *rsaCert) Generate(ky KeyType) error {
	var key *rsa.PrivateKey
	var err error
	switch ky {

	case RSA:
		rsa.
			key, err = rsa.GenerateKey(rand.Reader, 2048)
		c.data.Type = RSA256
	case RSA384:
		key, err = rsa.GenerateKey(rand.Reader, 3072)
		c.data.Type = RSA384
	case RSA512:
		key, err = rsa.GenerateKey(rand.Reader, 4096)
		c.data.Type = RSA512
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
	c.data.Key = marshal
	return nil
}

func (c *rsaCert) getPrivateKey() (*rsa.PrivateKey, error) {
	// Check if this is rsa type key
	switch c.data.Type {
	case RSA256:
	case RSA384:
	case RSA512:
	default:
		return nil, ErrWrongKeyType
	}

	key, err := x509.ParsePKCS8PrivateKey(c.data.Key)
	if err != nil {
		return nil, err
	}
	return key.(*rsa.PrivateKey), nil
}

func (c *rsaCert) GetPrivateKey() (interface{}, error) {
	return c.getPrivateKey()
}

func (c *rsaCert) getPublicKey() (*rsa.PublicKey, error) {
	privateKey, err := c.GetPrivateKey()
	if err != nil {
		return nil, err
	}
	return &privateKey.(*rsa.PrivateKey).PublicKey, nil
}

func (c *rsaCert) GetPublicKey() (interface{}, error) {
	return c.getPublicKey()
}

func (c *rsaCert) GetPublicStruct() (*PublicCertResponse, error) {
	publicKey, err := c.getPublicKey()
	if err != nil {
		return nil, err
	}
	marshaled, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return nil, err
	}
	thumbprint := sha1.Sum(marshaled)
	expBuf := new(bytes.Buffer)
	binary.Write(expBuf, binary.LittleEndian, int64(publicKey.E))
	pubStr := &PublicCertResponse{
		Alg: "RSA",
		Kty: string(c.data.Type),
		N:   base64.RawURLEncoding.EncodeToString(publicKey.N.Bytes()),
		E:   base64.RawURLEncoding.EncodeToString(expBuf.Bytes()),
		X5c: base64.StdEncoding.EncodeToString(marshaled),
		X5t: base64.RawURLEncoding.EncodeToString(thumbprint[:]),
		Kid: c.data.Id,
		Use: "sign",
	}
	return pubStr, nil
}

func (cr *rsaCert) GetSigningMethod() (jwt.SigningMethod, error) {
	method := jwt.GetSigningMethod(string(cr.data.AlgType))
	if method == nil {
		return nil, ErrWrongAlgType
	}
	return method, nil
}
