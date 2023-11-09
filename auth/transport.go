package auth

import (
	"net/http"

	"github.com/go-chi/chi"
	kithttp "github.com/go-kit/kit/transport/http"
)

func GetAuthHttpRouter(as AuthService) *chi.Mux {

	userAuthReq := kithttp.NewServer(
		makeUserAuthReq(as),
		decUserAuth,
		encUserAuth,
	)

	loginUserReq := kithttp.NewServer(
		makeLoginUserReq(as),
		decLoginUser,
		encLoginRes,
	)

	tokenRequest := kithttp.NewServer(
		makeTokenRequest(as),
		decServiceTokenReq,
		encTokenRes,
	)

	r := chi.NewRouter()
	r.Route("/authorize", func(r chi.Router) {
		// redirect to login page
		r.Get("/auth", userAuthReq.ServeHTTP)
		// Serve login page
		r.Get("/login", http.StripPrefix("/authorize/login", http.FileServer(http.Dir("html"))).ServeHTTP)
		// Post login data
		r.Post("/login", loginUserReq.ServeHTTP)
		// Get tokens
		r.Post("/token", tokenRequest.ServeHTTP)
	})
	return r
}
