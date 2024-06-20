package server

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

type ProxyResultCode string

const ProxyResultCookieOK = ProxyResultCode("CookieOK")
const ProxyResultBasicAuthOK = ProxyResultCode("BasicAuthOK")
const ProxyResultBasicAuthNG = ProxyResultCode("BasicAuthNG")
const ProxyResultFetchNG = ProxyResultCode("FetchNG")
const ProxyResultInternalError = ProxyResultCode("InternalError")

func (s *Server) proxyHandler(w http.ResponseWriter, r *http.Request) {
	resultCode := s.proxyMain(w, r)
	slog.Info("proxy done", "resultCode", resultCode)
}

func (s *Server) proxyMain(w http.ResponseWriter, r *http.Request) (resultCode ProxyResultCode) {
	// Get Cookie
	// Cookie Check
	cookieOk, err := s.Authenticater.IsValidCookie(r)
	if err != nil {
		slog.Error("get cookie", "error", err)
		return ProxyResultInternalError
	}

	if cookieOk {
		resultCode = ProxyResultCookieOK
	} else if s.Authenticater.CheckBasicAuth(r) { // BasicAuth Check
		resultCode = ProxyResultBasicAuthOK
	} else {
		// all NG
		w.Header().Add("WWW-Authenticate", `Basic realm="SECRET AREA"`)
		w.WriteHeader(http.StatusUnauthorized) // 401
		return ProxyResultBasicAuthNG
	}

	// To Proxy
	resp, err := s.Client.SendToProxy(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) // 500
		return ProxyResultFetchNG
	}
	defer func() {
		if resp.Body != nil {
			resp.Body.Close() // SendToProxy ではクローズしないのでここでクローズ
		}
	}()

	// Proxy Response ==> Server Response
	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))

	if resultCode != ProxyResultCookieOK { // Cookie OK 以外での認証OK出会った場合は Cookie 生成する
		// Generate Cookie
		cookie, err := s.Authenticater.GenerateCookie()
		if err != nil {
			return ProxyResultInternalError
		}
		http.SetCookie(w, cookie)
	}

	respBody := []byte("")
	if resp.Body != nil {
		respBody, err = io.ReadAll(resp.Body)
		if err != nil {
			return ProxyResultInternalError
		}
	}

	w.WriteHeader(resp.StatusCode)
	fmt.Fprint(w, string(respBody))
	return resultCode
}
