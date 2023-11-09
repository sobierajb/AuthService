package userRepo

import (
	"context"
	"smartHomeKit/common"

	"github.com/go-kit/kit/endpoint"
)

func makeCreateEndpoint(ur UserRepo) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UserPayload)
		user, err := ur.Create(&req)
		if err != nil {
			return common.Response[interface{}]{Data: nil, Error: err.Error()}, nil
		}
		return common.Response[*User]{Data: user, Error: ""}, nil
	}
}

func makeReadEndpoint(ur UserRepo) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(common.Id[string])
		user, err := ur.Read(req.Id)
		if err != nil {
			return common.Response[interface{}]{Data: nil, Error: err.Error()}, nil
		}
		return common.Response[*User]{Data: user, Error: ""}, nil
	}
}
func makeSearchEndpoint(ur UserRepo) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UserSearch)
		users, err := ur.Search(&req)
		if err != nil {
			return common.Response[interface{}]{Data: nil, Error: err.Error()}, nil
		}
		return common.Response[[]User]{Data: users, Error: ""}, nil
	}
}

func makeUpdateEndpoint(ur UserRepo) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(User)
		user, err := ur.Update(&req)
		if err != nil {
			return common.Response[interface{}]{Data: nil, Error: err.Error()}, nil
		}
		return common.Response[*User]{Data: user, Error: ""}, nil
	}
}

func makeDeleteEndpoint(ur UserRepo) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(common.Id[string])
		err := ur.Delete(req.Id)
		if err != nil {
			return common.Response[interface{}]{Data: nil, Error: err.Error()}, nil
		}
		return common.Response[string]{Data: "ok", Error: ""}, nil
	}
}
