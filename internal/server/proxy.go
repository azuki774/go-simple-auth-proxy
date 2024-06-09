package server

import (
	"azuki774/go-simple-auth-proxy/internal/auth"
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

func (s *Server) proxyHandler(w http.ResponseWriter, r *http.Request) {
	// Get Request

	// Authentication
	cookie := auth.GenerateCookie()

	// BasicAuth
	if !auth.CheckBasicAuth(r) {
		w.Header().Add("WWW-Authenticate", `Basic realm="SECRET AREA"`)
		w.WriteHeader(http.StatusUnauthorized) // 401
		return
	}

	// To Proxy
	resp, err := s.sendToProxy(r)
	defer resp.Body.Close()
	if err != nil {
		return
	}

	w.WriteHeader(resp.StatusCode)
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	http.SetCookie(w, cookie)
	fmt.Fprint(w, string(respBody))
}

// proxy先にリクエストを投げる。呼び出し元で resp を閉じること。
func (s *Server) sendToProxy(r *http.Request) (resp *http.Response, err error) {
	baseurl := r.URL.String()
	newurl := "http:" + "//" + s.ProxyAddr + baseurl
	slog.Info("newurl", "url", newurl)
	client := &http.Client{}

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(r.Method, newurl, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}
	// 呼び出し元で resp を閉じること: defer resp.Body.Close()
	return resp, nil
}
