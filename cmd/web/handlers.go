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

type deleteJobPostForm struct {
	Jpid string
}

// type JopPostFields struct {
// 	Position     string
// 	Description  string
// 	Contract     string
// 	Location     string
// 	Requirements struct {
// 		Content string
// 		Items   []string
// 	}
// 	Role struct {
// 		Content string
// 		Items   []string
// 	}
// 	validator.Validator
// }

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

	http.Redirect(w, r, "/user/signin", http.StatusSeeOther)

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
		// if user tries to create a new company with a name or a website of a company already created in db,
		// delete the user created in previous insert 🔝
		if errors.Is(err, models.ErrDuplicateCompanyName) || errors.Is(err, models.ErrDuplicateCompanyWebsite) {
			err := app.users.Delete(form.Email)
			if err != nil {
				app.serverError(w, r, err)
			}
		}
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

	http.Redirect(w, r, "/company/signin", http.StatusSeeOther)

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

	id, userType, err := app.users.Authenticate(form.Email, form.Password)
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
	app.sessionManager.Put(r.Context(), "userType", userType)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Remove(r.Context(), "authenticatedUserID")
	app.sessionManager.Remove(r.Context(), "userType")

	app.sessionManager.Put(r.Context(), "flash", "You've been logged out successfully.")

	http.Redirect(w, r, "/", http.StatusSeeOther)

}

func (app *application) userAccount(w http.ResponseWriter, r *http.Request) {

	// Get User info
	userId := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
	user, err := app.users.Get(userId)

	templateData := app.newTemplateData(r)
	templateData.User = user

	fmt.Println()
	fmt.Println("user Id", userId)
	fmt.Printf("user%+v", user)
	fmt.Println()

	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.Redirect(w, r, "/user/signin", http.StatusSeeOther)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	// Get either user job applications or user JobPosts

	if user.Type == models.UserTypeWorker {
		// todo:: get job applications, but first: create the job application page
	} else if user.Type == models.UserTypeCompany {
		// jobPosts := []models.JobPost{{ID: 1, Position: "iOS Engineer", Contract: "Full Time", Location: "United States", PostedAt: time.Now()}, {ID: 1, Position: "iOS Engineer", Contract: "Full Time", Location: "United States", PostedAt: time.Now()}, {ID: 1, Position: "iOS Engineer", Contract: "Full Time", Location: "United States", PostedAt: time.Now()}}
		fmt.Println("~~~~~~ get all company jobposts:::")
		jobPosts, err := app.users.GetJobPostByCompany(user.Name, 4)

		if err != nil {
			app.serverError(w, r, err)
		}

		fmt.Println()
		fmt.Printf("jobPosts: %+v\n", jobPosts)
		fmt.Println()

		templateData.JobPosts = jobPosts
	}

	app.render(w, r, http.StatusOK, "userAccount.tmpl.html", templateData)
}

func (app *application) userCreateJobPostGet(w http.ResponseWriter, r *http.Request) {
	// w.Write([]byte("render page for publish a jobpost"))

	data := app.newTemplateData(r)
	data.Form = models.CreateJopPostFields{}
	app.render(w, r, http.StatusOK, "createJobPost.tmpl.html", data)
}

func (app *application) userCreateJobPostPost(w http.ResponseWriter, r *http.Request) {
	var JP models.CreateJopPostFields

	err := decodeJSONBody(w, r, &JP)
	if err != nil {
		var mjr *malformedJSONRequest
		if errors.As(err, &mjr) {
			http.Error(w, mjr.msg, mjr.status)
		} else {
			app.serverError(w, r, err)
		}
	}

	fmt.Println()
	fmt.Printf("JP Fields: %+v\n", JP)
	fmt.Println()

	//> after this point we have the data available in JP variable, we need to validate that its present
	// {
	JP.CheckField(validator.NotBlank(JP.Position), "position", "This field can't be blank")
	JP.CheckField(validator.Matches(JP.Position, validator.LetterSpaceRegex), "position", "This field can only contain letters and spaces.")

	JP.CheckField(validator.NotBlank(JP.Description), "description", "This field can't be blank")
	JP.CheckField(validator.Matches(JP.Description, validator.LetterSpacesPunctuationRegex), "description", "This field can only contain letters, spaces and punctuation (, . ' ’ \" -)")

	JP.CheckField(validator.NotBlank(JP.Contract), "contract", "This field can't be blank")
	JP.CheckField(validator.PermittedValue(JP.Contract, "Full Time", "Part Time"), "contract", "This field must be either 'Full Time' or 'Part Time'")

	JP.CheckField(validator.NotBlank(JP.Location), "location", "This field can't be blank")

	JP.CheckField(validator.NotBlank(JP.Requirements.Content), "requirementsContent", "This field can't be blank")
	JP.CheckField(validator.Matches(JP.Requirements.Content, validator.LetterSpacesPunctuationExtendedNumbersRegex), "requirementsContent", "This field can only contain letters, numbers, spaces and punctuation (, . ' ’ \" & - ( ) /)")

	JP.CheckField(validator.ListHasItems(JP.Requirements.Items), "requirementsItems", "Please provide at least one item for the requirements list")

	for _, item := range JP.Requirements.Items {
		result := validator.Matches(item, validator.LetterSpacesPunctuationExtendedNumbersRegex)
		JP.CheckField(result, "requirementsItems", "Items for this list can only contain letters, numbers, spaces and punctuation (, . ' ’ \" & - ( ) /)")

		result = validator.Matches(item, validator.OnlyNumbersRegex)
		// fmt.Println("item", item, "result::", result)
		JP.CheckField(!result, "requirementsItems", "Items for this list cannot be only numbers.")

		result = validator.Matches(item, validator.OnlyPunctuationRegex)
		JP.CheckField(!result, "requirementsItems", "Items for this list cannot be only punctuation symbols.")
		if result {
			break
		}
	}

	JP.CheckField(validator.NotBlank(JP.Role.Content), "roleContent", "This field can't be blank")
	JP.CheckField(validator.Matches(JP.Role.Content, validator.LetterSpacesPunctuationExtendedNumbersRegex), "roleContent", "This field can only contain letters, numbers, spaces and punctuation (, . ' ’ \" & - ( ) /)")

	JP.CheckField(validator.ListHasItems(JP.Role.Items), "roleItems", "Please provide at least one item for the list of responsibilities")

	for _, item := range JP.Role.Items {
		result := validator.Matches(item, validator.LetterSpacesPunctuationExtendedNumbersRegex)
		JP.CheckField(result, "roleItems", "Items for this list can only contain letters, numbers, spaces and punctuation (, . ' ’ \" & - ( ) /)")

		result = validator.Matches(item, validator.OnlyNumbersRegex)
		JP.CheckField(!result, "roleItems", "Items for this list cannot be only numbers.")

		result = validator.Matches(item, validator.OnlyPunctuationRegex)
		JP.CheckField(!result, "roleItems", "Items for this list cannot be only punctuation symbols.")
		if result {
			break
		}
	}

	fmt.Println()
	fmt.Printf("JP fielderrors: %+v\n", JP.FieldErrors)
	fmt.Println()

	if !JP.Valid() {

		data := app.newTemplateData(r)
		data.Form = JP
		app.render(w, r, http.StatusUnprocessableEntity, "createJobPost.tmpl.html", data)
		return
	} else {
		//> after validation, data is ready to be inserted in the DB

		fmt.Printf("auth user id: %v\n", app.sessionManager.GetInt(r.Context(), "authenticatedUserID"))
		fmt.Printf("user type: %v\n", app.sessionManager.GetInt(r.Context(), "userType"))

	}
	// }

	userID := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")

	// {
	//TODO I have to put the data in the DB
	//* for that I can use the authenticatedUserID from user session, with that, query the -users_employers- table and get the -company_id- value
	//* then, with that value, check the -companies- table to see if there is a company for that ID, if there is one, we can proceed, but if not, we must return an error

	// with the data from user and the company id, I can to the following:
	//		insert a new record on requirements table (req_description, req_list)
	// 			get the id of that last record inserted in DB
	//		insert a new record on roles table (role_description, role_list)
	// 			get the id of that last record inserted in DB
	// put all data submitted by user, along with company_id and the last two values for req_id, and role_id

	// with that, it's been put on the DB and we can get it from company account page and the home page

	err = app.jobPosts.InsertJobPost(userID, JP)
	if err != nil {
		app.serverError(w, r, err)
	}
	app.sessionManager.Put(r.Context(), "flash", "Your JobPost has been published successfully.")

	http.Redirect(w, r, "/user/account", http.StatusSeeOther)

}

func (app *application) companyManageJobPosts(w http.ResponseWriter, r *http.Request) {

	// Get User info
	userId := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
	user, err := app.users.Get(userId)

	templateData := app.newTemplateData(r)
	templateData.User = user

	fmt.Println()
	fmt.Println("user Id", userId)
	fmt.Printf("user%+v", user)
	fmt.Println()

	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.Redirect(w, r, "/user/signin", http.StatusSeeOther)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	// jobPosts := []models.JobPost{{ID: 1, Position: "iOS Engineer", Contract: "Full Time", Location: "United States", PostedAt: time.Now()}, {ID: 1, Position: "iOS Engineer", Contract: "Full Time", Location: "United States", PostedAt: time.Now()}, {ID: 1, Position: "iOS Engineer", Contract: "Full Time", Location: "United States", PostedAt: time.Now()}}
	fmt.Println("~~~~~~ get all company jobposts:::")
	jobPosts, err := app.users.GetJobPostByCompany(user.Name, 0)

	if err != nil {
		app.serverError(w, r, err)
	}

	fmt.Println()
	fmt.Printf("jobPosts: %+v\n", jobPosts)
	fmt.Println()

	templateData.JobPosts = jobPosts

	app.render(w, r, http.StatusOK, "companyManageJobPosts.tmpl.html", templateData)
}

func (app *application) companyManageJobPostsDelete(w http.ResponseWriter, r *http.Request) {
	var form deleteJobPostForm

	err := decodeJSONBody(w, r, &form)
	if err != nil {
		var mjr *malformedJSONRequest
		if errors.As(err, &mjr) {
			http.Error(w, mjr.msg, mjr.status)
		} else {
			app.serverError(w, r, err)
		}
	}

	jpid, err := strconv.Atoi(form.Jpid)
	if err != nil {
		http.Error(w, "The jpid value must be an integer", http.StatusUnprocessableEntity)
	}

	fmt.Println()
	fmt.Printf("Job post to delete: jp with ID: %v", jpid)
	fmt.Println()

	err = app.jobPosts.DeleteJobPost(jpid)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Your publication has been deleted successfully")

	http.Redirect(w, r, "/account/manageJobPosts", http.StatusSeeOther)

}

func (app *application) userEditJobPost(w http.ResponseWriter, r *http.Request) {
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

	fmt.Println()
	fmt.Printf("%+v", jobPost)
	fmt.Println()

	data := app.newTemplateData(r)
	data.Form = models.CreateJopPostFields{}
	fmt.Println()
	fmt.Printf("%+v", data.Form)
	fmt.Println()
	data.JobPost = jobPost

	app.render(w, r, 200, "editJobPost.tmpl.html", data)
}

func (app *application) userEditJobPostPost(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	var JP models.EditJopPostFields

	err = decodeJSONBody(w, r, &JP)
	if err != nil {
		var mjr *malformedJSONRequest
		if errors.As(err, &mjr) {
			http.Error(w, mjr.msg, mjr.status)
		} else {
			app.serverError(w, r, err)
		}
	}

	fmt.Println()
	fmt.Printf("JP Fields: %+v\n", JP)
	fmt.Println()

	//> after this point we have the data available in JP variable, we need to validate that its present
	// {
	JP.CheckField(validator.NotBlank(JP.Position), "position", "This field can't be blank")
	JP.CheckField(validator.Matches(JP.Position, validator.LetterSpaceRegex), "position", "This field can only contain letters and spaces.")

	JP.CheckField(validator.NotBlank(JP.Description), "description", "This field can't be blank")
	JP.CheckField(validator.Matches(JP.Description, validator.LetterSpacesPunctuationRegex), "description", "This field can only contain letters, spaces and punctuation (, . ' ’ \" -)")

	JP.CheckField(validator.NotBlank(JP.Contract), "contract", "This field can't be blank")
	JP.CheckField(validator.PermittedValue(JP.Contract, "Full Time", "Part Time"), "contract", "This field must be either 'Full Time' or 'Part Time'")

	JP.CheckField(validator.NotBlank(JP.Location), "location", "This field can't be blank")

	JP.CheckField(validator.NotBlank(JP.Requirements.Content), "requirementsContent", "This field can't be blank")
	JP.CheckField(validator.Matches(JP.Requirements.Content, validator.LetterSpacesPunctuationExtendedNumbersRegex), "requirementsContent", "This field can only contain letters, numbers, spaces and punctuation (, . ' ’ \" & - ( ) /)")

	JP.CheckField(validator.ListHasItems(JP.Requirements.Items), "requirementsItems", "Please provide at least one item for the requirements list")

	for _, item := range JP.Requirements.Items {
		result := validator.Matches(item, validator.LetterSpacesPunctuationExtendedNumbersRegex)
		JP.CheckField(result, "requirementsItems", "Items for this list can only contain letters, numbers, spaces and punctuation (, . ' ’ \" & - ( ) /)")

		result = validator.Matches(item, validator.OnlyNumbersRegex)
		// fmt.Println("item", item, "result::", result)
		JP.CheckField(!result, "requirementsItems", "Items for this list cannot be only numbers.")

		result = validator.Matches(item, validator.OnlyPunctuationRegex)
		JP.CheckField(!result, "requirementsItems", "Items for this list cannot be only punctuation symbols.")
		if result {
			break
		}
	}

	JP.CheckField(validator.NotBlank(JP.Role.Content), "roleContent", "This field can't be blank")
	JP.CheckField(validator.Matches(JP.Role.Content, validator.LetterSpacesPunctuationExtendedNumbersRegex), "roleContent", "This field can only contain letters, numbers, spaces and punctuation (, . ' ’ \" & - ( ) /)")

	JP.CheckField(validator.ListHasItems(JP.Role.Items), "roleItems", "Please provide at least one item for the list of responsibilities")

	for _, item := range JP.Role.Items {
		result := validator.Matches(item, validator.LetterSpacesPunctuationExtendedNumbersRegex)
		JP.CheckField(result, "roleItems", "Items for this list can only contain letters, numbers, spaces and punctuation (, . ' ’ \" & - ( ) /)")

		result = validator.Matches(item, validator.OnlyNumbersRegex)
		JP.CheckField(!result, "roleItems", "Items for this list cannot be only numbers.")

		result = validator.Matches(item, validator.OnlyPunctuationRegex)
		JP.CheckField(!result, "roleItems", "Items for this list cannot be only punctuation symbols.")
		if result {
			break
		}
	}

	fmt.Println()
	fmt.Printf("JP fielderrors: %+v\n", JP.FieldErrors)
	fmt.Println()

	if !JP.Valid() {

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
		data.Form = JP
		data.JobPost = jobPost

		// ?????????????????????????????????????????????????????????????????????????????????????????????
		// TODO: get the info for the JobPost field in data
		app.render(w, r, http.StatusUnprocessableEntity, "editJobPost.tmpl.html", data)
		return
	}
	//> after validation, data is ready to be inserted in the DB

	fmt.Printf("auth user id: %v\n", app.sessionManager.GetInt(r.Context(), "authenticatedUserID"))
	fmt.Printf("user type: %v\n", app.sessionManager.GetInt(r.Context(), "userType"))
	// }

	// userID := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")

	// {
	//TODO I have to put the data in the DB
	//* for that I can use the authenticatedUserID from user session, with that, query the -users_employers- table and get the -company_id- value
	//* then, with that value, check the -companies- table to see if there is a company for that ID, if there is one, we can proceed, but if not, we must return an error

	// with the data from user and the company id, I can to the following:
	//		insert a new record on requirements table (req_description, req_list)
	// 			get the id of that last record inserted in DB
	//		insert a new record on roles table (role_description, role_list)
	// 			get the id of that last record inserted in DB
	// put all data submitted by user, along with company_id and the last two values for req_id, and role_id

	// with that, it's been put on the DB and we can get it from company account page and the home page

	// err = app.jobPosts.InsertJobPost(userID, JP)
	//??? err = app.jobPosts.EditJobPost(JP.ID, JP.CompanyID, JP.Requirements.ReqID, JP.Role.RoleID, JP)
	err = app.jobPosts.EditJobPost(JP)
	if err != nil {
		app.serverError(w, r, err)
	}
	app.sessionManager.Put(r.Context(), "flash", "Your JobPost has been updated successfully.")

	http.Redirect(w, r, "/user/account", http.StatusSeeOther)

}
