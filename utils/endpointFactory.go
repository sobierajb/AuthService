package utils

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type ServiceFunction[T interface{}, G interface{}] func(T) (G, error)

func EndpointFactory[T interface{}, G interface{}](sfunc ServiceFunction[T, G]) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(T)
		res, err := sfunc(req)
		if err != nil {
			return nil, err
		}
		return res, nil
	}
}
