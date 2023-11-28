package clientRepo

import (
	"smartHomeKit/repository/certRepo"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type GrantType string

const (
	GtAutorisationCode  GrantType = "autorisation_code"
	GtRefreshToken      GrantType = "refresh_token"
	GtClientCredentials GrantType = "client_credentials"
)

type Client struct {
	Id              string
	Name            string
	GrantType       GrantType
	RedirectUri     string
	Secret          string
	Certificate     certRepo.Certificate // If null then use random certificate stored in repository
	AccessTokenExp  time.Duration
	RefreshTokenExp time.Duration
}

type ClientRepo interface {
	Create(*Client) (*Client, error)
	Update(*Client) (*Client, error)
	GetById(id string) (*Client, error)
	Delete(id string)
}

func (c *Client) GetAccessTokenExp(claim jwt.MapClaims) jwt.MapClaims {
	if claim == nil {
		claim = jwt.MapClaims{}
	}
	claim["iat"] = time.Now().Unix()
	claim["exp"] = time.Now().Add(c.AccessTokenExp).Unix()
	return claim
}

func (c *Client) GetRefreshTokenExp(claim jwt.MapClaims) jwt.MapClaims {
	if claim == nil {
		claim = jwt.MapClaims{}
	}
	claim["iat"] = time.Now().Unix()
	claim["exp"] = time.Now().Add(c.AccessTokenExp).Unix()
	return claim
}
