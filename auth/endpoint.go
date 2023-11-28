package auth

import (
	"context"
	"smartHomeKit/utils"

	"github.com/go-kit/kit/endpoint"
)

func makeUserAuthReq(as AuthService) endpoint.Endpoint {
	return utils.EndpointFactory[UserAuthRequest, UserAuthRequest](as.UserAuthReq)
}

func makeLoginUserReq(as AuthService) endpoint.Endpoint {
	return utils.EndpointFactory[LoginRequest, UserAuthResponse](as.LoginUserReq)
}

func makeTokenRequest(as AuthService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		switch t := request.(type) {
		case TokenRequest[AuthCodeRequest]:
			res, err := as.AuthCodeReq(*t.Data)
			if err != nil {
				return nil, err
			}
			return res, nil
		case TokenRequest[ClientCredentialsRequest]:
			res, err := as.ClientCredsReq(*t.Data)
			if err != nil {
				return nil, err
			}
			return res, nil
		case TokenRequest[RefreshTokenRequest]:
			res, err := as.RefreshTokenReq(*t.Data)
			if err != nil {
				return nil, err
			}
			return res, nil
		default:
			return nil, ErrUndefinedTokenRequestPayload
		}
	}
}
