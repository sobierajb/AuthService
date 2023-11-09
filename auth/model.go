package auth

import (
	"net/url"

	"github.com/golang-jwt/jwt"
)

const (
	grAuthorisationCode = "authorisation_code"
	grRefreshToken      = "refresh_token"
	grClientCredentials = "client_credentials"
)

type UserAuthRequest struct {
	Challenge   string
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

type UserTokenRequest struct {
	ClientId     string
	CodeVerifier string
	Code         string
}

type ServiceTokenRequest struct {
	ClientId     string
	ClientSecret string
}

type TokenResponse struct {
	AccessToken  jwt.Token
	RefreshToken jwt.Token
	Expires      int
}

type LoginRequest struct {
	Login    string
	Password string
	*UserAuthRequest
}
