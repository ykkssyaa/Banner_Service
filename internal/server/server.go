package server

import (
	"BannerService/internal/service"
	"net/http"
	"time"
)

type HttpServer struct {
	services *service.Services
	//logger   *logger.Logger
}

func NewHttpServer(addr string, services *service.Services) *http.Server {
	server := &HttpServer{services: services}

	r := NewRouter(server)

	return &http.Server{
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
		Addr:           addr,
	}
}
