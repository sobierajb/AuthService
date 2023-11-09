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
		var res TokenResponse
		var err error
		switch req := request.(type) {
		case TokenRequest[UserTokenRequest]:
			res, err = as.UserTokeReq(*req.Data)
			if err != nil {
				return nil, err
			}
			return res, nil
		case TokenRequest[ServiceTokenRequest]:
			res, err = as.ServiceTokenReq(*req.Data)
			if err != nil {
				return nil, err
			}
			return res, nil
		default:
			return nil, ErrUnknownTokenReqType
		}
	}
}

// func makeUserTokenRequest(as AuthService) endpoint.Endpoint {
// 	return utils.EndpointFactory[UserTokenRequest, TokenResposne](as.UserTokeReq)
// }

// func makeServiceAuthRequest(as AuthService) endpoint.Endpoint {
// 	return utils.EndpointFactory[ServiceAuthRequest, TokenResposne](as.ServiceAuthReq)
// }
