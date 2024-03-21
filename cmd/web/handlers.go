package main

import "net/http"

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	//w.Write([]byte("main page, here i'll put the list of job posts"))

	jobposts, err := app.jobPosts.Latest()

	if err != nil {
		app.serverError(w, r, err)
		return
	}

	templateData := app.newTemplateData()
	templateData.JobPosts = jobposts
	app.render(w, r, 200, "home.tmpl.html", templateData)
}
