package server

import (
	"net/http"
)

type mockClient struct {
	err error
}
type mockAuthenticater struct {
	basicok  bool
	cookieok bool
}

type mockResponseWriter struct {
}

func (m *mockClient) SendToProxy(r *http.Request) (resp *http.Response, err error) {
	if m.err != nil {
		return nil, m.err
	}

	return &http.Response{
		Status:     "200 OK",
		StatusCode: 200,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Body:       nil, // この関数では読まれないので nil でいいはず
	}, nil
}

func (m *mockAuthenticater) GenerateCookie() (*http.Cookie, error) {
	cookie := &http.Cookie{
		Name:  "token",
		Value: "example_token_value", // TODO
	}
	return cookie, nil
}

func (m *mockAuthenticater) IsValidCookie(r *http.Request) (ok bool, err error) {
	return m.cookieok, nil
}

func (m *mockAuthenticater) CheckBasicAuth(r *http.Request) bool {
	return m.basicok
}

func (m *mockResponseWriter) Header() http.Header {
	return http.Header{}
}

func (m *mockResponseWriter) Write([]byte) (int, error) {
	return 1, nil
}

func (m *mockResponseWriter) WriteHeader(statusCode int) {
}
