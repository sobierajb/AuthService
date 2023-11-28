package auth

import (
	"crypto/rand"
	"encoding/base64"
	"smartHomeKit/repository/certRepo"
	"smartHomeKit/repository/clientRepo"
	"smartHomeKit/repository/sessionRepo"
	"smartHomeKit/repository/userRepo"

	"github.com/golang-jwt/jwt/v4"
)

type AuthService interface {
	// redirect to login page
	UserAuthReq(UserAuthRequest) (UserAuthRequest, error)
	LoginUserReq(LoginRequest) (UserAuthResponse, error)
	AuthCodeReq(AuthCodeRequest) (TokenResponse, error)
	ClientCredsReq(ClientCredentialsRequest) (TokenResponse, error)
	RefreshTokenReq(RefreshTokenRequest) (TokenResponse, error)
}

type authService struct {
	userRepo    userRepo.UserRepo
	sessionRepo sessionRepo.SessionRepo
	certRepo    certRepo.CertRepo
	clientRepo  clientRepo.ClientRepo
}

// Constructor for auth service
func NewAuthService(ur userRepo.UserRepo, sr sessionRepo.SessionRepo, cr certRepo.CertRepo, cl clientRepo.ClientRepo) *authService {
	return &authService{
		userRepo:    ur,
		sessionRepo: sr,
		certRepo:    cr,
		clientRepo:  cl,
	}
}

// First step in authorisation code flow - redirect to proper login page
func (as *authService) UserAuthReq(uar *UserAuthRequest) (*UserAuthRequest, error) {
	// Check if client is registerd in DB
	_, err := as.clientRepo.GetById(uar.ClientId)
	if err != nil {
		return nil, err
	}
	return uar, nil
}

// Second step in authorisation code flow. Check if user exists. Compare password hashes.
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
		Client:        &clientRepo.Client{Id: lr.ClientId},
		User:          user,
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

// Third step in authorisation code flow. Check if code is correct.
func (as *authService) AuthCodeReq(acr *AuthCodeRequest) (*TokenResponse, error) {

	//Calidate if given code was registered as session
	session, err := as.sessionRepo.GetSessionByCode(acr.Code)
	if err != nil {
		return nil, err
	}
	//validate code verifier. If verifier not valid, return err
	if err := session.ValidateCodeChalange(acr.CodeVerifier); err != nil {
		return nil, err
	}

	method, err := session.Client.Certificate.GetSigningMethod()
	if err != nil {
		return nil, err
	}

	// Create accessToken and all required claims
	accessToken := jwt.New(method)
	claims := session.Client.GetAccessTokenExp(jwt.MapClaims{})
	session.User.GetUserClaims(claims)
	accessToken.Claims = claims

	refreshToken := jwt.New(method)
	refreshToken.Header["kid"] = session.Client.Certificate.GetThumbprint()
	refreshToken.Claims = session.Client.GetRefreshTokenExp(jwt.MapClaims{})
	key, err := session.Client.Certificate.GetPrivateKey()
	if err != nil {
		return nil, err
	}
	signedAccessToken, err := accessToken.SignedString(key)
	if err != nil {
		return nil, err
	}
	signedRefreshToke, err := refreshToken.SignedString(key)
	if err != nil {
		return nil, err
	}
	tokenResponse := &TokenResponse{
		AccessToken:  signedAccessToken,
		RefreshToken: signedRefreshToke,
	}
	return tokenResponse, nil
}

func (as *authService) ClientCredsReq(str *ClientCredentialsRequest) (*TokenResponse, error) {
	client, err := as.clientRepo.GetById(str.ClientId)
	if err != nil {
		return nil, err
	}

	if client.Secret != str.ClientSecret {
		return nil, ErrCannotAuthorize
	}
	method, err := client.Certificate.GetSigningMethod()
	accessToken := jwt.New(method)
	accessToken.Claims = client.GetAccessTokenExp(jwt.MapClaims{})

	refreshToken := jwt.New(method)
	refreshToken.Claims = client.GetRefreshTokenExp(jwt.MapClaims{})

	key, err := client.Certificate.GetPrivateKey()
	if err != nil {
		return nil, err
	}
	signedAccessToken, err := accessToken.SignedString(key)
	if err != nil {
		return nil, err
	}

	signedRefreshToken, err := refreshToken.SignedString(key)
	if err != nil {
		return nil, err
	}

	tokenResponse := &TokenResponse{
		AccessToken:  signedAccessToken,
		RefreshToken: signedRefreshToken,
	}
	return tokenResponse, nil
}

// Function for generating random code by given length of bytes. Return base64 Raw URL Encoding
func genRandomCode(length int) (string, error) {
	data := make([]byte, length)
	n, err := rand.Read(data)
	if n != length || err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(data), nil
}

// Refresh token authorisation flow
func (as *authService) RefreshTokenReq(rfr *RefreshTokenRequest) (*TokenResponse, error) {
	client, err := as.clientRepo.GetById(rfr.ClientId)
	if err != nil {
		return nil, err
	}

	// Verify and parse refresh token
	token, err := jwt.Parse(rfr.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != string(client.Certificate.GetAlghorythm()) {
			return nil, ErrWrongSigningMethod
		}
		key, err := client.Certificate.GetPublicKey()
		if err != nil {
			return nil, err
		}
		return &key, nil
	})

	if err != nil {
		return nil, err
	}
	// Read claims
	claims := token.Claims.(jwt.RegisteredClaims)
	// get user
	user, _ := as.userRepo.GetById(claims.Subject)

	method, err := client.Certificate.GetSigningMethod()
	if err != nil {
		return nil, err
	}
	key, err := client.Certificate.GetPrivateKey()
	if err != nil {
		return nil, err
	}

	accessToken := jwt.New(method)
	accessToken.Claims = client.GetAccessTokenExp(jwt.MapClaims{})
	// If user was present in refresh token then add user claims informations
	if user != nil {
		accessToken.Claims = user.GetUserClaims(accessToken.Claims.(jwt.MapClaims))
	}
	signedAccessToken, err := accessToken.SignedString(key)
	if err != nil {
		return nil, err
	}

	refreshToken := jwt.New(method)
	refreshToken.Claims = client.GetRefreshTokenExp(jwt.MapClaims{})
	signedRefreshToken, err := refreshToken.SignedString(key)
	if err != nil {
		return nil, err
	}
	tokenResponse := &TokenResponse{
		AccessToken:  signedAccessToken,
		RefreshToken: signedRefreshToken,
	}
	return tokenResponse, nil
}
