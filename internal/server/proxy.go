package server

import (
	"azuki774/go-simple-auth-proxy/internal/auth"
	"fmt"
	"net/http"
)

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	// Get Request

	// Authentication
	cookie := auth.GenerateCookie()

	// BasicAuth
	if !auth.CheckBasicAuth(r) {
		w.Header().Add("WWW-Authenticate", `Basic realm="SECRET AREA"`)
		w.WriteHeader(http.StatusUnauthorized) // 401
		return
	}

	// Response
	http.SetCookie(w, cookie)
	fmt.Fprintf(w, "Hello, World\n")
}
