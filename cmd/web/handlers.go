package main

import (
	"fmt"
	"net/http"
)

type jobPostFilterForm struct {
	Position string `form:"position"`
	Location string `form:"location"`
	Contract string `form:"contract"`
}

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

func (app *application) homeFilterJobPosts(w http.ResponseWriter, r *http.Request) {

	var form jobPostFilterForm
	err := app.decodeForm(w, r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	fmt.Printf("%+v\n", form)
}
