package utils

import (
	"context"
	"encoding/json"
	"net/http"

	kithttp "github.com/go-kit/kit/transport/http"
)

type RequestDecoder func(context.Context, *http.Request) (interface{}, error)
type ResponseEncoder func(context.Context, http.ResponseWriter, interface{}) error

func GetRequestDecoder[T interface{}]() kithttp.DecodeRequestFunc {

	var request T
	return func(ctx context.Context, r *http.Request) (interface{}, error) {
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			return nil, err
		}
		return request, nil
	}
}

func GetResponseEncoder() ResponseEncoder {
	return func(ctx context.Context, w http.ResponseWriter, response interface{}) error {
		return json.NewEncoder(w).Encode(response)
	}
}

func EncodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
