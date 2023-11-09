package userRepo

import (
	"context"
	"encoding/json"
	"net/http"
	"smartHomeKit/utils"

	"github.com/go-chi/chi"
	kithttp "github.com/go-kit/kit/transport/http"
)

func MakeHttpUserHandlers(ur UserRepo) *chi.Mux {

	read := kithttp.NewServer(
		makeReadEndpoint(ur),
		readDecoder,
		utils.EncodeResponse,
	)

	search := kithttp.NewServer(
		makeSearchEndpoint(ur),
		searchDecoder,
		utils.EncodeResponse,
	)

	create := kithttp.NewServer(
		makeCreateEndpoint(ur),
		createDecoder,
		utils.EncodeResponse,
	)

	update := kithttp.NewServer(
		makeUpdateEndpoint(ur),
		updateDecoder,
		utils.EncodeResponse,
	)

	delete := kithttp.NewServer(
		makeDeleteEndpoint(ur),
		deleteDecoder,
		utils.EncodeResponse,
	)

	r := chi.NewRouter()
	r.Route("/users", func(r chi.Router) {
		r.Get("/", read.ServeHTTP)
		r.Get("/search", search.ServeHTTP)
		r.Post("/", create.ServeHTTP)
		r.Put("/{userId}", update.ServeHTTP)
		r.Delete("/{userId}", delete.ServeHTTP)
	})
	return r
}

func readDecoder(ctx context.Context, r *http.Request) (interface{}, error) {
	userId := chi.URLParam(r, "id")
	if userId == "" {
		return nil, ErrEmptyParam
	}
	return userId, nil
}

func searchDecoder(ctx context.Context, r *http.Request) (interface{}, error) {
	var userSearch UserSearch
	if err := json.NewDecoder(r.Body).Decode(&userSearch); err != nil {
		return nil, err
	}
	return userSearch, nil

}

func createDecoder(ctx context.Context, r *http.Request) (interface{}, error) {
	var userPayload UserPayload
	if err := json.NewDecoder(r.Body).Decode(&userPayload); err != nil {
		return nil, err
	}
	return userPayload, nil
}

func updateDecoder(ctx context.Context, r *http.Request) (interface{}, error) {
	var user User
	user.Id = chi.URLParam(r, "userId")

	if user.Id == "" {
		return nil, ErrEmptyParam
	}

	if err := json.NewDecoder(r.Body).Decode(&user.UserPayload); err != nil {
		return nil, err
	}
	return user, nil
}

func deleteDecoder(ctx context.Context, r *http.Request) (interface{}, error) {
	userId := chi.URLParam(r, "userId")
	if userId == "" {
		return nil, ErrEmptyParam
	}
	return userId, nil
}
