package main

import (
	"net/http"

	"github.com/go-chi/chi"
)

func main() {

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hi server"))
	}

	r := chi.NewRouter()
	r.Get("/", handler)
	fileServer := http.FileServer(http.Dir("./html/"))
	r.Handle("/authorize/*", http.StripPrefix("/authorize/", fileServer))

	// r.Mount("/", auth.GetAuthHttpRouter())
	http.ListenAndServe(":8080", r)

}
