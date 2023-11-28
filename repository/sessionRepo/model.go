package sessionRepo

import (
	"bytes"
	"crypto"
	"encoding/base64"
	"smartHomeKit/repository/clientRepo"
	"smartHomeKit/repository/userRepo"
	"time"
)

type ChallengeMethod string

const (
	S256Method  ChallengeMethod = "S256"
	PlainMethod ChallengeMethod = "plain"
)

type Session struct {
	Id              string
	RemoteAddr      string
	Client          *clientRepo.Client
	User            *userRepo.User
	Code            string
	CodeChallenge   string
	ChallengeMethod ChallengeMethod
	State           string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (s *Session) ValidateCodeChalange(verifier string) error {
	var decodedChallenge []byte

	// decode code challenge to byte array
	if _, err := base64.RawURLEncoding.Decode(decodedChallenge, []byte(s.CodeChallenge)); err != nil {
		return ErrCannotValidateCodeChallenge
	}

	switch s.ChallengeMethod {
	case S256Method:

		hashedVerifier := crypto.SHA256.New()
		hashedVerifier.Write([]byte(verifier))
		if bytes.Equal(decodedChallenge, hashedVerifier.Sum(nil)) {
			return nil
		}
	case PlainMethod:
		// code verifier = code chalenge
		if s.CodeChallenge == verifier {
			return nil
		}
	default:
		return ErrUknownChallengeMethod
	}
	return ErrCannotValidateCodeChallenge
}
