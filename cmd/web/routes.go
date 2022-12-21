package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()
	fileServer := http.FileServer(neuteredFS{http.Dir("./ui/static")})

	mux.Handle("/static/",
		http.StripPrefix("/static", fileServer))
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	middleWare := alice.New(
		app.recoverPanic,
		app.logRequest,
		secureHeaders,
	)
	return middleWare.Then(mux)
}
