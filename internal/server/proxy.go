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

	// Response
	http.SetCookie(w, cookie)
	fmt.Fprintf(w, "Hello, World\n")
}
