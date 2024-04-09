package main

import (
	"BannerService/internal/server"
	lg "BannerService/pkg/logger"
	"errors"
	"net/http"
)

func main() {
	logger := lg.InitLogger()

	logger.Info.Print("Creating server.")

	// port := viper.GetString("PORT")
	port := "8080"
	srv := server.NewHttpServer(":" + port)

	logger.Info.Print("Starting the server on port: " + port + "\n\n")
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Err.Fatalf("error occured while running http server: \"%s\" \n", err.Error())
	}
}
