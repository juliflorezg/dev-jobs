package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/juliflorezg/dev-jobs/internal/models"
	"github.com/juliflorezg/dev-jobs/internal/validator"
)

type JobPostFilterForm struct {
	Position            string `form:"position"`
	Location            string `form:"location"`
	Contract            string `form:"contract"`
	validator.Validator `form:"-"`
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
	templateData.Form = JobPostFilterForm{}
	app.render(w, r, 200, "home.tmpl.html", templateData)
}

func (app *application) homeFilterJobPosts(w http.ResponseWriter, r *http.Request) {

	var form JobPostFilterForm
	err := app.decodeForm(w, r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	fmt.Printf("%+v\n", form)

	form.ValidateFormData((validator.IsNoData(form.Position, form.Location, form.Contract)))

	fmt.Println()
	fmt.Println("is form valid?? ::", form.Valid())
	fmt.Printf("form ::>%+v \n", form)
	fmt.Println()

	if !form.Valid() {

		jobposts, err := app.jobPosts.Latest()

		if err != nil {
			app.serverError(w, r, err)
			return
		}

		data := app.newTemplateData()
		data.JobPosts = jobposts

		data.Form = form

		fmt.Println()
		fmt.Printf("form when error ::>%+v \n", form)
		fmt.Println()

		app.render(w, r, http.StatusUnprocessableEntity, "home.tmpl.html", data)

		return
	}

	templateData := app.newTemplateData()
	templateData.Form = form
	jobPosts, err := app.jobPosts.FilterPosts(form.Position, form.Location, form.Contract)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			templateData.JobPostsFilterData.NoPostsData = "There are not matching results for your search criteria"
		} else {

			app.serverError(w, r, err)
			return
		}
	}

	templateData.JobPosts = jobPosts
	templateData.JobPostsFilterData.IsSearchResultPage = true

	msg := getSearchResultMessage(form.Position, form.Location, form.Contract)
	templateData.JobPostsFilterData.SearchResultMessage = msg
	app.render(w, r, 200, "home.tmpl.html", templateData)

}
