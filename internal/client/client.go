package client

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"time"
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

	// set meta to proxy
	req.Header.Set("Content-Type", r.Header.Get("Content-Type"))

	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// URI: GET / にアクセスして疎通があることを確認する
// ステータスコードは問わない
func (c *Client) Ping() (err error) {
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", c.ProxyAddr, nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer func() {
		io.Copy(io.Discard, resp.Body) // 読み捨てる
		resp.Body.Close()
	}()

	return nil
}
