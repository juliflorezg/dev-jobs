package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/form/v4"
)

type malformedJSONRequest struct {
	msg    string
	status int
}

func (mjr *malformedJSONRequest) Error() string {
	return mjr.msg
}

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
		err := fmt.Errorf("the template %s does not exist", page)
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

func (app *application) newTemplateData(r *http.Request) templateData {

	return templateData{
		CurrentYear:     time.Now().Year(),
		Flash:           app.sessionManager.PopString(r.Context(), "flash"),
		IsAuthenticated: app.isAuthenticated(r),
		UserType:        app.getUserType(r),
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
		//handles case where user doesn't upload a file, later in handlers svg icon will be set to a default value
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

type HSL struct {
	Hue        int
	Saturation int
	Lightness  int
}

func HexToHSL(hex string) HSL {
	r, _ := strconv.ParseInt(hex[1:3], 16, 0)
	g, _ := strconv.ParseInt(hex[3:5], 16, 0)
	b, _ := strconv.ParseInt(hex[5:7], 16, 0)

	var rVal, gVal, bVal float64
	rVal = float64(r) / 255
	gVal = float64(g) / 255
	bVal = float64(b) / 255

	max := max(rVal, gVal, bVal)
	min := min(rVal, gVal, bVal)

	var h, s, l float64
	h = 0
	s = 0
	l = (max + min) / 2

	if max == min {
		h = 0
		s = 0
	} else {
		d := max - min
		if l > 0.5 {
			s = d / (2 - max - min)
		} else {
			s = d / (max + min)
		}

		switch max {
		case rVal:
			h = (gVal - bVal) / d
			if gVal < bVal {
				h += 6
			}
		case gVal:
			h = (bVal-rVal)/d + 2
		case bVal:
			h = (rVal-gVal)/d + 4
		}
		h /= 6
	}

	return HSL{
		Hue:        int(h * 360),
		Saturation: int(s * 100),
		Lightness:  int(l * 100),
	}
}

func max(a, b, c float64) float64 {
	max := a
	if b > max {
		max = b
	}
	if c > max {
		max = c
	}
	return max
}

func min(a, b, c float64) float64 {
	min := a
	if b < min {
		min = b
	}
	if c < min {
		min = c
	}
	return min
}

func GetHSLColorStr(hsl HSL) string {
	return fmt.Sprintf("hsl(%v, %v%%, %v%%)", hsl.Hue, hsl.Saturation, hsl.Lightness)
}

func formatCompanyName(name string) string {
	nameLower := strings.ToLower(strings.Trim(name, " "))
	nameWords := strings.Split(nameLower, " ")
	formattedName := ""
	for i, word := range nameWords {
		if i < len(word)-1 {
			formattedName += strings.ToUpper(string(word[0])) + word[1:] + " "
		} else {
			formattedName += strings.ToUpper(string(word[0])) + word[1:]
		}
	}

	return strings.Trim(formattedName, " ")
}

func (app *application) isAuthenticated(r *http.Request) bool {
	return app.sessionManager.Exists(r.Context(), "authenticatedUserID")
}

func (app *application) getUserType(r *http.Request) int {
	return app.sessionManager.GetInt(r.Context(), "userType")
}

func decodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) error {

	//check for content type header
	ct := r.Header.Get("Content-Type")
	fmt.Println()
	fmt.Println(ct)
	fmt.Println()
	if ct != "" {
		mediaType := strings.ToLower(strings.TrimSpace(strings.Split(ct, ";")[0]))
		if mediaType != "application/json" {
			msg := "Content-Type header is not application/json"
			return &malformedJSONRequest{status: http.StatusUnsupportedMediaType, msg: msg}
		}
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			return &malformedJSONRequest{status: http.StatusBadRequest, msg: msg}

		case errors.Is(err, io.ErrUnexpectedEOF):
			msg := "Request body contains badly-formed JSON"
			return &malformedJSONRequest{status: http.StatusBadRequest, msg: msg}

		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			return &malformedJSONRequest{status: http.StatusBadRequest, msg: msg}

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
			return &malformedJSONRequest{status: http.StatusBadRequest, msg: msg}

		case errors.Is(err, io.EOF):
			msg := "Request body must not be empty"
			return &malformedJSONRequest{status: http.StatusBadRequest, msg: msg}

		case err.Error() == "http: request body too large":
			msg := "Request body must not be larger than 1MB"
			return &malformedJSONRequest{status: http.StatusBadRequest, msg: msg}

		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{})
	if !errors.Is(err, io.EOF) {
		msg := "Request body must only contain a single JSON object"
		return &malformedJSONRequest{status: http.StatusBadRequest, msg: msg}
	}

	return nil
}
