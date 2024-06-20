package client

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
)

type Client struct {
	ProxyAddr string
}

// proxy先にリクエストを投げる。
// r は元クライアントからもらったリクエスト
// resp.Body は 呼び出し元でクローズすること
func (c *Client) SendToProxy(r *http.Request) (resp *http.Response, err error) {
	baseurl := r.URL.String()
	newurl := c.ProxyAddr + baseurl
	slog.Info("proxy to", "url", newurl)
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

	return resp, nil
}
