package server

import (
	"BannerService/internal/service"
	"net/http"
	"time"
)

type HttpServer struct {
	services *service.Services
}

func NewHttpServer(addr string) *http.Server {
	server := &HttpServer{}

	r := NewRouter(server)

	return &http.Server{
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
		Addr:           addr,
	}
}
