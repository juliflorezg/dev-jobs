package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/go-playground/form/v4"
)

func ReadUserIP(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}
	return IPAddress
}

// The clientError helper sends a specific status code and corresponding description
// to the user. We'll use this to send responses like 400 "Bad
// Request" when there's a problem with the request that the user sent
func (app *application) clientError(w http.ResponseWriter, statusCode int) {
	http.Error(w, http.StatusText(statusCode), statusCode)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
		trace  = string(debug.Stack())
	)

	app.logger.Error(err.Error(), slog.String("method", method), slog.String("uri", uri), slog.String("trace", trace))

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) render(w http.ResponseWriter, r *http.Request, status int, page string, data templateData) {
	ts, ok := app.templateCache[page]

	if !ok {
		err := fmt.Errorf("The template %s does not exist", page)
		app.serverError(w, r, err)
		return
	}

	buf := new(bytes.Buffer)
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	w.WriteHeader(status)
	buf.WriteTo(w)
}

func (app *application) newTemplateData() templateData {

	return templateData{
		CurrentYear: time.Now().Year(),
	}
}

// Create a new decodePostForm() helper method. The second parameter here, dst,
// is the target destination that we want to decode the form data into.
func (app *application) decodeForm(w http.ResponseWriter, r *http.Request, dst any) error {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return err
	}

	// Call Decode() on our decoder instance, passing the target destination as
	// the first parameter

	fmt.Println("in decodeForm fn, helpers.go")
	fmt.Printf("+%v\n", r.Form)
	fmt.Printf("+%v\n", r.PostForm)
	fmt.Println("in helpers file, decodeForm fn (end)")

	if method := r.Method; method == "GET" {
		err = app.formDecoder.Decode(dst, r.Form)
	} else {
		err = app.formDecoder.Decode(dst, r.PostForm)
	}

	if err != nil {
		// If we try to use an invalid target destination, the Decode() method
		// will return an error with the type *form.InvalidDecoderError.
		// We use errors.As() to check for this and raise a panic rather than returning
		// the error.

		var invalidDecoderError *form.InvalidDecoderError

		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}

		// For all other errors, we return them as normal.
		return err
	}

	return nil
}

func processFile(r *http.Request) (string, error) {
	var s string
	f, h, err := r.FormFile("svgicon")
	if err != nil {
		return s, err
	}
	defer f.Close()

	// for your information
	fmt.Println("\nfile:", f, "\nheader:", h, "\nerr", err)

	// read
	bs, err := io.ReadAll(f)
	if err != nil {
		return s, err
	}
	s = string(bs)

	return s, nil
}

func getSearchResultMessage(position, location, contract string) []string {

	msg := []string{"Search results for"}
	if position == "" {
		msg = append(msg, "position:", "all open positions")
	} else {
		msg = append(msg, "position:", position)
	}
	if location == "" {
		msg = append(msg, "location:", "all countries")
	} else {
		msg = append(msg, "location:", location)
	}
	if contract == "" {
		msg = append(msg, "contract:", "full time & part time")
	} else {
		msg = append(msg, "contract:", contract)
	}

	return msg
}
