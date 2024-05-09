package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/juliflorezg/dev-jobs/ui"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	fileServer := http.FileServer(http.FS(ui.Files))
	router.Handler(http.MethodGet, "/static/*filepath", fileServer)

	dynamic := alice.New(app.sessionManager.LoadAndSave)

	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/filterBy", dynamic.ThenFunc(app.homeFilterJobPosts))
	router.Handler(http.MethodGet, "/jobpost/view/:id", dynamic.ThenFunc(app.jobPostView))

	router.Handler(http.MethodGet, "/user/signup", dynamic.ThenFunc(app.userSignUp))
	router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(app.userSignUpPost))
	router.Handler(http.MethodGet, "/company/signup", dynamic.ThenFunc(app.companySignUp))
	router.Handler(http.MethodPost, "/company/signup", dynamic.ThenFunc(app.companySignUpPost))

	router.Handler(http.MethodGet, "/user/signin", dynamic.ThenFunc(app.userLogin))
	router.Handler(http.MethodGet, "/company/signin", dynamic.ThenFunc(app.companyLogin))

	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return standard.Then(router)
}
