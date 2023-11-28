package sessionRepo

import "errors"

var (
	ErrUknownChallengeMethod       = errors.New("unknown chalenge method")
	ErrSessionNotFound             = errors.New("session not found")
	ErrCannotValidateCodeChallenge = errors.New("cannot validate code challange")
)
