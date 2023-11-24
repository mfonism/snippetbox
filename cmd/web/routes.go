package main

import (
	"net/http"

	"github.com/mfonism/snippetbox/ui"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	// custom handler for 404 Not Found responses
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	fileServer := http.FileServer(neuteredFileSystem{http.FS(ui.Files)})
	router.Handler(http.MethodGet, "/static/*filepath", fileServer)

	// middleware chains
	dynamicMiddlewareChain := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate)
	requireAuthenticationMiddlewareChain := dynamicMiddlewareChain.Append(app.requireAuthentication)
	makeRequireLogoutMiddlewareChain := func(action string) alice.Chain {
		return dynamicMiddlewareChain.Append(app.requireLogout(action))
	}

	router.Handler(
		http.MethodGet,
		"/",
		dynamicMiddlewareChain.ThenFunc(app.home),
	)
	router.Handler(
		http.MethodGet,
		"/snippet/view/:id",
		dynamicMiddlewareChain.ThenFunc(app.snippetView),
	)

	// require logout
	router.Handler(
		http.MethodGet,
		"/user/signup",
		makeRequireLogoutMiddlewareChain("sign up").ThenFunc(app.userSignup),
	)
	router.Handler(
		http.MethodPost,
		"/user/signup",
		makeRequireLogoutMiddlewareChain("sign up").ThenFunc(app.userSignupPost),
	)
	router.Handler(
		http.MethodGet,
		"/user/login",
		makeRequireLogoutMiddlewareChain("log in").ThenFunc(app.userLogin),
	)
	router.Handler(
		http.MethodPost,
		"/user/login",
		makeRequireLogoutMiddlewareChain("log in").ThenFunc(app.userLoginPost),
	)

	// require authentication
	router.Handler(
		http.MethodGet,
		"/snippet/create",
		requireAuthenticationMiddlewareChain.ThenFunc(app.snippetCreate),
	)
	router.Handler(
		http.MethodPost,
		"/snippet/create",
		requireAuthenticationMiddlewareChain.ThenFunc(app.snippetCreatePost),
	)
	router.Handler(
		http.MethodPost,
		"/user/logout",
		requireAuthenticationMiddlewareChain.ThenFunc(app.userLogoutPost),
	)

	standardMiddlewareChain := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return standardMiddlewareChain.Then(router)
}
