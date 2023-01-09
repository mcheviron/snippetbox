package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	filerServer := http.FileServer(http.Dir("./ui/static"))
	router.Handler(
		http.MethodGet, "/static/*filepath",
		app.noDirListingHandler(http.StripPrefix("/static", filerServer)),
	)

	// This will be a middleware added to our main routes to create a per user session
	dynamicMiddleware := alice.New(
		app.sessionManager.LoadAndSave,
		noSurf,
		app.authenticate,
	)
	router.Handler(
		http.MethodGet,
		"/",
		dynamicMiddleware.ThenFunc(app.home),
	)
	router.Handler(
		http.MethodGet,
		"/snippet/view/:id",
		dynamicMiddleware.ThenFunc(app.snippetView),
	)
	router.Handler(
		http.MethodGet,
		"/user/signup",
		dynamicMiddleware.ThenFunc(app.userSignup),
	)
	router.Handler(
		http.MethodPost,
		"/user/signup",
		dynamicMiddleware.ThenFunc(app.userSignupPost),
	)
	router.Handler(
		http.MethodGet,
		"/user/login",
		dynamicMiddleware.ThenFunc(app.userLogin),
	)
	router.Handler(
		http.MethodPost,
		"/user/login",
		dynamicMiddleware.ThenFunc(app.userLoginPost),
	)

	// NOTE: Protected (authenticated-only) routes
	protectedMiddleware := dynamicMiddleware.Append(app.requireAuthentication)
	router.Handler(
		http.MethodGet,
		"/snippet/create",
		protectedMiddleware.ThenFunc(app.snippetCreate),
	)
	router.Handler(
		http.MethodPost,
		"/snippet/create",
		protectedMiddleware.ThenFunc(app.snippetCreatePost),
	)
	router.Handler(
		http.MethodPost,
		"/user/logout",
		protectedMiddleware.ThenFunc(app.userLogoutPost),
	)
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	return standardMiddleware.Then(router)
}
