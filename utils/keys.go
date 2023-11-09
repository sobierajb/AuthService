package utils

type ctxKey int

const (
	NoValue ctxKey = iota
	StoredRequest
)
