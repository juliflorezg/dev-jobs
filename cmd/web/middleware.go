package main

import "net/http"

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var (
			ip     = ReadUserIP(r)
			proto  = r.Proto
			method = r.Method
			uri    = r.URL.RequestURI()
		)

		app.logger.Info("Received a new request ->", "ip", ip, "protocol", proto, "method", method, "uri", uri)

		next.ServeHTTP(w, r)
	})
}
