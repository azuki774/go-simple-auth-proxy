package server

import (
	"azuki774/go-simple-auth-proxy/internal/metrics"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type Client interface {
	SendToProxy(r *http.Request) (resp *http.Response, err error)
	Ping() (err error)
}

type Authenticater interface {
	GenerateCookie() (*http.Cookie, error)
	IsValidCookie(r *http.Request) (ok bool, err error)
	CheckBasicAuth(r *http.Request) bool
}
type Server struct {
	ListenPort    string
	Client        Client // ex. http://example.com:1234
	Authenticater Authenticater
	ExporterPort  string
}

func (s *Server) Start(ctx context.Context) (err error) {
	slog.Info("server start")

	go func() { // prometheus exporter はエラーハンドリングしない
		if s.ExporterPort == "" {
			return
		}
		m := metrics.NewMetrics(s.ExporterPort)
		slog.Info("start metrics server")
		_ = m.Start()
		slog.Info("close metrics server")
	}()

	addr := fmt.Sprintf(":%s", s.ListenPort)
	http.HandleFunc("/", s.proxyHandler) // ハンドラを登録してウェブページを表示させる

	server := &http.Server{
		Addr:    addr,
		Handler: nil,
	}

	ctxIn, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	var errCh = make(chan error)
	go func() {
		errCh <- server.ListenAndServe()
	}()

	<-ctxIn.Done()
	if nerr := server.Shutdown(ctx); nerr != nil {
		slog.Error("failed to shutdown server", "err", err)
		return nerr
	}

	err = <-errCh
	if err != nil && err != http.ErrServerClosed {
		slog.Error("failed to close server", "err", err)
		return err
	}

	slog.Info("http server close gracefully")
	return nil
}

// Server が start しても問題ないコンフィグ設定になっているかを確認する
func (s *Server) CheckReadiness() (err error) {
	// proxy 先の正常性確認
	if err := s.Client.Ping(); err != nil {
		slog.Error("failed to proxyAddr", "err", err)
		return err
	}
	return nil
}
