package certRepo

import "errors"

var (
	ErrWrongKeyType = errors.New("wrong key type")
	ErrCertNotFound = errors.New("certificate with given id not found")
	ErrWrongAlgType = errors.New("wrong alghorytm type")
)
