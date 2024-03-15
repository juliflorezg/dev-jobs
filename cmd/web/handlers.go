package main

import "net/http"

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	//w.Write([]byte("main page, here i'll put the list of job posts"))

	templateData := app.newTemplateData()
	app.render(w, r, 200, "home.tmpl.html", templateData)
}
