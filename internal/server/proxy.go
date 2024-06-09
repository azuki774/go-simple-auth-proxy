package server

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

type ProxyResultCode string

const ProxyResultOK = ProxyResultCode("OK")
const ProxyResultBasicAuthNG = ProxyResultCode("BasicAuthNG")
const ProxyResultFetchNG = ProxyResultCode("FetchNG")
const ProxyResultInternalError = ProxyResultCode("InternalError")

func (s *Server) proxyHandler(w http.ResponseWriter, r *http.Request) {
	resultCode := s.proxyMain(w, r)
	slog.Info("proxy done", "resultCode", resultCode)
}

func (s *Server) proxyMain(w http.ResponseWriter, r *http.Request) (resultCode ProxyResultCode) {
	// Get Cookie

	// BasicAuth
	if !s.Authenticater.CheckBasicAuth(r) {
		w.Header().Add("WWW-Authenticate", `Basic realm="SECRET AREA"`)
		w.WriteHeader(http.StatusUnauthorized) // 401
		return ProxyResultBasicAuthNG
	}

	// To Proxy
	resp, err := s.Client.SendToProxy(r)
	defer resp.Body.Close() // SendToProxy ではクローズしないのでここでクローズ
	if err != nil {
		return ProxyResultFetchNG
	}

	// Proxy Response ==> Server Response
	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	w.WriteHeader(resp.StatusCode)

	// Generate Cookie
	cookie := s.Authenticater.GenerateCookie()
	http.SetCookie(w, cookie)

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return ProxyResultInternalError
	}

	fmt.Fprint(w, string(respBody))
	return ProxyResultOK
}
