package auth

import "errors"

var ErrParseLoginUrl = errors.New("cannot parse login url")
var ErrNoRequestCtx = errors.New("cannot get http request pointer from context")
var ErrUnknownTokenReqType = errors.New("unknown token request type")
var ErrUnknownGrantType = errors.New("unknown grant type")
var ErrClientNotFound = errors.New("client with given id not registered")
var ErrCannotAuthorize = errors.New("cannot authorize")
var ErrCannotGenCode = errors.New("cannot generate authorisation code")
var ErrWrongSigningMethod = errors.New("token signing method does not match")
var ErrUndefinedTokenRequestPayload = errors.New("undefined token request payload")
