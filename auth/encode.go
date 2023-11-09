package auth

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"smartHomeKit/utils"
)

var redirectUrl = "SomeLoginUrlWhichHasToBest"

// Redirect To login Page
func encUserAuth(ctx context.Context, w http.ResponseWriter, response interface{}) error {

	req, ok := ctx.Value(utils.StoredRequest).(*http.Request)
	if !ok {
		return ErrNoRequestCtx
	}
	res := response.(UserAuthRequest)
	u, err := url.Parse(redirectUrl)
	if err != nil {
		return errors.Join(err, ErrParseLoginUrl)
	}
	u.Query().Set("challange", res.Challenge)
	u.Query().Set("client_id", res.ClientId)
	u.Query().Set("redirect_uri", res.RedirectUri.String())
	u.Query().Set("state", res.State)
	http.Redirect(w, req, u.String(), http.StatusMovedPermanently)
	return nil
}

func encLoginRes(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	req, ok := ctx.Value(utils.StoredRequest).(*http.Request)
	if !ok {
		return ErrNoRequestCtx
	}
	res := response.(UserAuthResponse)
	http.Redirect(w, req, res.RedirectUri.String(), http.StatusFound)
	return nil
}

func encTokenRes(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
