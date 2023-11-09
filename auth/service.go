package auth

import (
	"crypto/rand"
	"encoding/base64"
	"smartHomeKit/repository/sessionRepo"
	"smartHomeKit/repository/userRepo"
)

type AuthService interface {
	// redirect to login page
	UserAuthReq(UserAuthRequest) (UserAuthRequest, error)
	LoginUserReq(LoginRequest) (UserAuthResponse, error)
	UserTokeReq(UserTokenRequest) (TokenResponse, error)
	ServiceTokenReq(ServiceTokenRequest) (TokenResponse, error)
}

type authService struct {
	userRepo    userRepo.UserRepo
	sessionRepo sessionRepo.SessionRepo
}

func NewAuthService(ur userRepo.UserRepo, sr sessionRepo.SessionRepo) *authService {
	return &authService{
		userRepo:    ur,
		sessionRepo: sr,
	}
}

func (as *authService) UserAuthReq(uar *UserAuthRequest) (*UserAuthRequest, error) {
	// Check if client is registerd in DB
	_, err := as.sessionRepo.GetClientById(uar.ClientId)
	if err != nil {
		return nil, err
	}
	return uar, nil
}

func (as *authService) LoginUserReq(lr *LoginRequest) (*UserAuthResponse, error) {
	// Find User in DB
	user, err := as.userRepo.GetByLogin(lr.Login)
	if err != nil {
		return nil, err
	}
	// Check if password is matched
	err = user.CheckPassword(lr.Password)
	if err != nil {
		return nil, err
	}

	// Generate Code
	code, err := genRandomCode(64)
	if err != nil {
		return nil, ErrCannotGenCode
	}
	// Create Session
	var session = &sessionRepo.Session{
		ClientId:      lr.ClientId,
		Code:          code,
		CodeChallenge: lr.Challenge,
		State:         lr.State,
		RemoteAddr:    lr.RemoteAddr,
	}

	session, err = as.sessionRepo.CreateSession(session)
	if err != nil {
		return nil, err
	}

	// create and return Resposne
	return &UserAuthResponse{
		Code:        session.Code,
		RedirectUri: *lr.RedirectUri,
		State:       lr.State,
	}, nil
}

func (as *authService) UserTokenRequest(utr *UserTokenRequest) (*TokenResponse, error) {

	token := &TokenResponse{}
	return token, nil
}

func (as *authService) ServiceTokenRequest(str *ServiceTokenRequest) (*TokenResponse, error) {
	client, err := as.sessionRepo.GetClientById(str.ClientId)
	if err != nil {
		return nil, err
	}
	if client.Secret != str.ClientSecret {
		return nil, ErrCannotAuthorize
	}

	token := &TokenResponse{}
	return token, nil
}

func genRandomCode(length int) (string, error) {
	data := make([]byte, length)
	n, err := rand.Read(data)
	if n != length || err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data), nil
}
