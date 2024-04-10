package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/juliflorezg/dev-jobs/internal/models"
	"github.com/juliflorezg/dev-jobs/internal/validator"
)

type JobPostFilterForm struct {
	Position            string `form:"position"`
	Location            string `form:"location"`
	Contract            string `form:"contract"`
	validator.Validator `form:"-"`
}

type userSignUpForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
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

func (app *application) jobPostView(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	jobPost, err := app.jobPosts.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	data := app.newTemplateData()
	data.JobPost = jobPost
	app.render(w, r, 200, "viewJobPost.tmpl.html", data)

}

func (app *application) userSignUp(w http.ResponseWriter, r *http.Request) {
	// w.Write([]byte("sample response"))
	data := app.newTemplateData()
	data.Form = userSignUpForm{}
	app.render(w, r, http.StatusOK, "userSignUp.tmpl.html", data)
}
func (app *application) userSignUpPost(w http.ResponseWriter, r *http.Request) {

	var form userSignUpForm

	err := app.decodeForm(w, r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRegex), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characters")
	form.CheckField(validator.Matches(form.Password, validator.SpecialCharsRegex), "password", "Your password must contain at least: one uppercase letter, one lowercase letter, one number and one special character (!\"@#$%^&*()?<>.-)")

	if !form.Valid() {
		data := app.newTemplateData()
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "userSignUp.tmpl.html", data)
		return
	}

	err = app.users.Insert(form.Name, form.Email, form.Password, 1)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "This email address is already in use")
			data := app.newTemplateData()
			data.Form = form
			app.render(w, r, http.StatusUnprocessableEntity, "userSignUp.tmpl.html", data)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Your sign up was successful. Please log in.")

  http.Redirect(w, r, "/user/login", http.StatusSeeOther)

	// w.Write([]byte("create a new user..."))

}
func (app *application) companySignUp(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("sample response"))
}
func (app *application) companySignUpPost(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("sample response"))
}
