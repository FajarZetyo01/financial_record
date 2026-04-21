package helpers

import (
	"financial_record/config"
	"net/http"
)

func GuestOnly(next http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		loogedIn := config.SessionManager.GetBool(request.Context(), "LOGGED_IN")
		if loogedIn {
			http.Redirect(writer, request, "/home", http.StatusSeeOther)
		}
		next.ServeHTTP(writer, request)
	}
}

func AuthOnly(next http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		loogedIn := config.SessionManager.GetBool(request.Context(), "LOGGED_IN")
		if !loogedIn {
			http.Redirect(writer, request, "/login", http.StatusSeeOther)
		}
		next.ServeHTTP(writer, request)
	}
}
