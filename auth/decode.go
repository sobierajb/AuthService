package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
)

func parseUserAuthRequest(r *http.Request) (*UserAuthRequest, error) {
	values := r.URL.Query()
	// store pointer to request in context to
	redirectUri, err := url.Parse(values.Get("redirect_uri"))
	if err != nil {
		return nil, err
	}

	userAuthReq := &UserAuthRequest{
		Challenge:   values.Get("challange"),
		ClientId:    values.Get("client_id"),
		RedirectUri: redirectUri,
		State:       values.Get("state"),
		RemoteAddr:  r.RemoteAddr,
	}
	return userAuthReq, nil
}

func decUserAuth(ctx context.Context, r *http.Request) (interface{}, error) {
	return parseUserAuthRequest(r)
}

func decLoginUser(ctx context.Context, r *http.Request) (interface{}, error) {
	formValues := r.PostForm
	userRequest, err := parseUserAuthRequest(r)
	if err != nil {
		return nil, err
	}

	loginReq := &LoginRequest{
		Login:           formValues.Get("username"),
		Password:        formValues.Get("password"),
		UserAuthRequest: userRequest,
	}
	return loginReq, nil
}

func decTokenRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	formValues := r.PostForm
	grantType := formValues.Get("grant_type")
	switch grantType {
	case grAuthorisationCode:
		// crate tokenRequest object
		tr := &TokenRequest[UserTokenRequest]{
			GrantType: grantType,
			Data: &UserTokenRequest{
				ClientId:     formValues.Get("client_id"),
				CodeVerifier: formValues.Get("code_verifier"),
				Code:         formValues.Get("code"),
			},
		}
		return tr, nil
	case grClientCredentials:
		tr := &TokenRequest[ServiceTokenRequest]{
			GrantType: grantType,
			Data: &ServiceTokenRequest{
				ClientId:     formValues.Get("client_id"),
				ClientSecret: formValues.Get("client_secret"),
			},
		}
		return tr, nil
	case grRefreshToken:
		return "", nil
	default:
		return nil, ErrUnknownGrantType
	}
}

func decUserTokenReq(ctx context.Context, r *http.Request) (interface{}, error) {
	var req UserTokenRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func decServiceTokenReq(ctx context.Context, r *http.Request) (interface{}, error) {
	var req ServiceTokenRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}
