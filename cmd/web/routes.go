package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/juliflorezg/dev-jobs/ui"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	fileServer := http.FileServer(http.FS(ui.Files))
	router.Handler(http.MethodGet, "/static/*filepath", fileServer)

	router.HandlerFunc(http.MethodGet, "/", app.home)
	router.HandlerFunc(http.MethodGet, "/filterBy", app.homeFilterJobPosts)
	router.HandlerFunc(http.MethodGet, "/jobpost/view/:id", app.jobPostView)

	return app.recoverPanic(app.logRequest(secureHeaders(router)))
}
