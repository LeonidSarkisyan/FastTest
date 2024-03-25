package server

import (
	"context"
	"crypto/tls"
	"net/http"
	"time"
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run(port string, handler http.Handler) error {
	tlsConfig := &tls.Config{}

	cert, err := tls.LoadX509KeyPair("/etc/letsencrypt/live/xn----7sbq0ccdbck.xn--p1ai/fullchain.pem", "/etc/letsencrypt/live/xn----7sbq0ccdbck.xn--p1ai/privkey.pem")
	if err != nil {

		s.httpServer = &http.Server{
			Addr:           ":80",
			Handler:        handler,
			MaxHeaderBytes: 1 << 20,
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
		}

		return s.httpServer.ListenAndServe()

	}
	tlsConfig.Certificates = []tls.Certificate{cert}

	s.httpServer = &http.Server{
		Addr:           port,
		Handler:        handler,
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		TLSConfig:      tlsConfig,
	}

	return s.httpServer.ListenAndServeTLS("", "")
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
