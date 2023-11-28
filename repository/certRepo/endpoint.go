package certRepo

import (
	"context"
	"smartHomeKit/common"

	"github.com/go-kit/kit/endpoint"
)

func makeCreateEndpoint(cp CertRepo) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(UserPayload)
		user, err := ur.Create(&req)
		if err != nil {
			return common.Response[interface{}]{Data: nil, Error: err.Error()}, nil
		}
		return common.Response[*User]{Data: user, Error: ""}, nil
	}
}
