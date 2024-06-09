package server

import (
	"fmt"
	"io"
	"net/http"
)

func (s *Server) proxyHandler(w http.ResponseWriter, r *http.Request) {
	// Get Request

	// Generate Cookie
	cookie := s.Authenticater.GenerateCookie()

	// BasicAuth
	if !s.Authenticater.CheckBasicAuth(r) {
		w.Header().Add("WWW-Authenticate", `Basic realm="SECRET AREA"`)
		w.WriteHeader(http.StatusUnauthorized) // 401
		return
	}

	// To Proxy
	resp, err := s.Client.SendToProxy(r)
	defer resp.Body.Close() // SendToProxy ではクローズしないのでここでクローズ
	if err != nil {
		return
	}

	// Proxy Response ==> Server Response
	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	w.WriteHeader(resp.StatusCode)
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	http.SetCookie(w, cookie)
	fmt.Fprint(w, string(respBody))
}
