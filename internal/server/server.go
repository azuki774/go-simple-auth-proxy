package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type Server struct {
	ListenPort string
	MyAddr     string // http://localhost:8888
	ProxyAddr  string // MyAddr -> ProxyAddr に置換
}

func (s *Server) Start(ctx context.Context) (err error) {
	slog.Info("server start")

	// TODO:
	s.MyAddr = "localhost:8080"
	s.ProxyAddr = "localhost:8888"
	///

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
