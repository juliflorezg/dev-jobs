package main

import (
	"errors"
	"fmt"
	"log"
	"math"
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
type companySignUpForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	SVGIcon             string `form:"svgicon"`
	IconBgColor         string `form:"iconbgcolor"`
	Website             string `form:"companywebsite"`
	validator.Validator `form:"-"`
}

type userLoginForm struct {
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

	templateData := app.newTemplateData(r)
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

		data := app.newTemplateData(r)
		data.JobPosts = jobposts

		data.Form = form

		fmt.Println()
		fmt.Printf("form when error ::>%+v \n", form)
		fmt.Println()

		app.render(w, r, http.StatusUnprocessableEntity, "home.tmpl.html", data)

		return
	}

	templateData := app.newTemplateData(r)
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

	data := app.newTemplateData(r)
	data.JobPost = jobPost
	app.render(w, r, 200, "viewJobPost.tmpl.html", data)

}

func (app *application) userSignUp(w http.ResponseWriter, r *http.Request) {
	// w.Write([]byte("sample response"))
	data := app.newTemplateData(r)
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
	form.CheckField(validator.IsValidPassword(form.Password), "password", "Your password must contain at least: one uppercase letter, one lowercase letter, one number and one special character (!\"@#$%^&*()?<>.-)")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "userSignUp.tmpl.html", data)
		return
	}

	err = app.users.Insert(form.Name, form.Email, form.Password, models.UserTypeWorker)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "This email address is already in use")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, r, http.StatusUnprocessableEntity, "userSignUp.tmpl.html", data)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Your sign-up was successful. Please log in.")

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)

	// w.Write([]byte("create a new user..."))

}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = companySignUpForm{}
	app.render(w, r, http.StatusOK, "userLogin.tmpl.html", data)
}

func (app *application) companySignUp(w http.ResponseWriter, r *http.Request) {
	// w.Write([]byte("sample response"))
	data := app.newTemplateData(r)
	data.Form = companySignUpForm{}
	app.render(w, r, http.StatusOK, "companySignUp.tmpl.html", data)
}

func (app *application) companySignUpPost(w http.ResponseWriter, r *http.Request) {
	// w.Write([]byte("sample response"))
	var form companySignUpForm

	// err := app.decodeForm(w, r, &form)
	// if err != nil {
	// 	app.clientError(w, http.StatusBadRequest)
	// 	return
	// }

	s, err := processFile(r)
	if err != nil {
		// app.serverError(w, r, err)
		s = models.DefaultCompanyIcon
	}

	r.ParseMultipartForm(math.MaxInt16)

	log.Printf("%#v", r.Header.Get("Content-Type"))

	name := r.PostFormValue("name")
	email := r.PostFormValue("email")
	passwd := r.PostFormValue("password")
	iconBgColor := r.PostFormValue("iconbgcolor")
	website := r.PostFormValue("companywebsite")
	form = companySignUpForm{
		Name:        name,
		Email:       email,
		Password:    passwd,
		SVGIcon:     s,
		IconBgColor: GetHSLColorStr(HexToHSL(iconBgColor)),
		Website:     website,
	}

	// log.Println(name)
	// log.Println(email)
	// log.Println(passwd)
	// log.Println(iconBgColor)
	// log.Println(website)

	// fmt.Println()
	// fmt.Printf("%+v", form)
	// fmt.Println()

	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRegex), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characters")
	form.CheckField(validator.IsValidPassword(form.Password), "password", "Your password must contain at least: one uppercase letter, one lowercase letter, one number and one special character (!\"@#$%^&*()?<>.-)")
	form.CheckField(validator.Matches(form.Website, validator.WebsiteRegex), "website", "Your company website doesn't match a URL pattern, try one of this formats: http://example.com, https://example.com, http://example.com/xyz, https://example.com.xyz http://www.example.com, https://www.example.com, http://www.example.com/xyz, https://www.example.com/xyz")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "companySignUp.tmpl.html", data)
		return
	}

	form.Name = formatCompanyName(form.Name)

	// fmt.Println("form NAME:::", form.Name)

	err = app.users.Insert(form.Name, form.Email, form.Password, models.UserTypeCompany)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "This email address is already in use")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, r, http.StatusUnprocessableEntity, "companySignUp.tmpl.html", data)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	err = app.users.InsertCompany(form.Name, form.SVGIcon, form.IconBgColor, form.Website)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateCompanyName) {
			fmt.Println(err)
			form.AddFieldError("name", "There is already a company with this name.")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, r, http.StatusUnprocessableEntity, "companySignUp.tmpl.html", data)
		} else if errors.Is(err, models.ErrDuplicateCompanyWebsite) {
			fmt.Println(err)
			form.AddFieldError("website", "There is already a company with this website.")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, r, http.StatusUnprocessableEntity, "companySignUp.tmpl.html", data)
		}
		app.serverError(w, r, err)
		return
	}

	// if user and company where created successfully we can insert a record in users_employers table with users.id and companies.company_id
	usrId, compId, err := app.users.GetLastUserCompanyCreated(form.Email, form.Name)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	err = app.users.InsertCompanyUser(usrId, compId)
	if err != nil {
		app.serverError(w, r, err)
	}

	app.sessionManager.Put(r.Context(), "flash", "Your sign-up was successful Please log in.")

	http.Redirect(w, r, "/company/login", http.StatusSeeOther)

	// io.WriteString(w, s+name+email+passwd+"-|-"+iconBgColor+"-|-"+website)
}

func (app *application) companyLogin(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = companySignUpForm{}
	app.render(w, r, http.StatusOK, "companyLogin.tmpl.html", data)
}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	var form userLoginForm

	err := app.decodeForm(w, r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Email), "email", "Please provide your email to sign in.")
	form.CheckField(validator.Matches(form.Email, validator.EmailRegex), "email", "This field must be a valid email address.")
	form.CheckField(validator.NotBlank(form.Password), "password", "Please provide your password to sign in.")

	id, err := app.users.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Your email or password is incorrect")

			data := app.newTemplateData(r)
			data.Form = form

			// fmt.Println(id)
			if r.URL.RequestURI() == "/user/signin" {
				app.render(w, r, http.StatusOK, "userLogin.tmpl.html", data)
			} else if r.URL.RequestURI() == "/company/signin" {
				app.render(w, r, http.StatusOK, "companyLogin.tmpl.html", data)
			}
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "authenticatedUserID", id)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
