package server

import (
	"net/http"
	"time"
)

type HTTPServer struct {
	server http.Server
}

func (s *HTTPServer) Run(port string, handler http.Handler) error {
	s.server = http.Server{
		Addr:              ":" + port,
		Handler:           handler,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 0,
		WriteTimeout:      10 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}

	return s.server.ListenAndServe()
}
