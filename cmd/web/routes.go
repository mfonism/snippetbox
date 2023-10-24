package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	// custom handler for 404 Not Found responses
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	fileServer := http.FileServer(neuteredFileSystem{http.Dir("./ui/static/")})
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	dynamicMiddlewareChain := alice.New(app.sessionManager.LoadAndSave)

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
	router.Handler(
		http.MethodGet,
		"/snippet/create",
		dynamicMiddlewareChain.ThenFunc(app.snippetCreate),
	)
	router.Handler(
		http.MethodPost,
		"/snippet/create",
		dynamicMiddlewareChain.ThenFunc(app.snippetCreatePost),
	)

	router.Handler(
		http.MethodGet,
		"/user/signup",
		dynamicMiddlewareChain.ThenFunc(app.userSignup),
	)
	router.Handler(
		http.MethodPost,
		"/user/signup",
		dynamicMiddlewareChain.ThenFunc(app.userSignupPost),
	)
	router.Handler(
		http.MethodGet,
		"/user/login",
		dynamicMiddlewareChain.ThenFunc(app.userLogin),
	)
	router.Handler(
		http.MethodPost,
		"/user/login",
		dynamicMiddlewareChain.ThenFunc(app.userLoginPost),
	)
	router.Handler(
		http.MethodPost,
		"/user/logout",
		dynamicMiddlewareChain.ThenFunc(app.userLogoutPost),
	)

	standardMiddlewareChain := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return standardMiddlewareChain.Then(router)
}
