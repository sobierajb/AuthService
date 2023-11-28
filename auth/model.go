package auth

import (
	"net/url"
)

const (
	grAuthorisationCode = "authorisation_code"
	grRefreshToken      = "refresh_token"
	grClientCredentials = "client_credentials"
)

type UserAuthRequest struct {
	Challenge   string
	Method      string
	ClientId    string
	RedirectUri *url.URL
	State       string
	RemoteAddr  string
}

type UserAuthResponse struct {
	RedirectUri url.URL
	Code        string
	State       string
}

type TokenRequest[T interface{}] struct {
	GrantType string
	Data      *T
}

type AuthCodeRequest struct {
	ClientId     string
	CodeVerifier string
	Code         string
}

type ClientCredentialsRequest struct {
	ClientId     string
	ClientSecret string
}

type RefreshTokenRequest struct {
	ClientId     string
	RefreshToken string
}

type LoginRequest struct {
	Login    string
	Password string
	*UserAuthRequest
}

type TokenResponse struct {
	AccessToken  string
	RefreshToken string
}
