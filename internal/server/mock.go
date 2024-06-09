package server

import (
	"net/http"
)

type mockClient struct {
}
type mockAuthenticater struct {
}

type mockResponseWriter struct {
}

func (m *mockClient) SendToProxy(r *http.Request) (resp *http.Response, err error) {
	return &http.Response{
		Status:     "200 OK",
		StatusCode: 200,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Body:       nil, // この関数では読まれないので nil でいいはず
	}, nil
}

func (m *mockAuthenticater) GenerateCookie() *http.Cookie {
	cookie := &http.Cookie{
		Name:  "token",
		Value: "example_token_value", // TODO
	}
	return cookie
}

func (m *mockAuthenticater) IsValidCookie(r *http.Request) (ok bool, err error) {
	return true, nil
}

func (m *mockAuthenticater) CheckBasicAuth(r *http.Request) bool {
	return true
}

func (m *mockResponseWriter) Header() http.Header {
	return http.Header{}
}

func (m *mockResponseWriter) Write([]byte) (int, error) {
	return 1, nil
}

func (m *mockResponseWriter) WriteHeader(statusCode int) {
}
