package certRepo

import "time"

type AlghorythmType string

const (
	RSA256   AlghorythmType = "RS256"
	RSA384   AlghorythmType = "RS384"
	RSA512   AlghorythmType = "RS512"
	HMAC256  AlghorythmType = "HS256"
	HMAC384  AlghorythmType = "HS384"
	HMAC512  AlghorythmType = "HS512"
	ECDSA256 AlghorythmType = "ES256"
	ECDSA384 AlghorythmType = "ES384"
	ECDSA512 AlghorythmType = "ES512"
	EdDSA    AlghorythmType = "EdDSA"
)

type CertificateData struct {
	Id         string
	AlgType    AlghorythmType
	Key        []byte
	Thumbprint string
	CreatedAt  time.Time
}

type certifcatePayload struct {
	AlghorythmType
}

type PublicCertResponse struct {
	Alg string `json:"alg"` // cryptographic alghorythm used with key
	Kty string `json:"kty"` // familly of cryptographic alghorytm
	N   string `json:"n"`   //modulus for RSA public
	E   string `json:"e"`   //exponent for RSA public
	X5t string `json:"x5t"` //The thumbprint
	X5c string `json:"x5c"` //x.509 chain
	Use string `json:"use"`
	Kid string `json:"kid"`
}
