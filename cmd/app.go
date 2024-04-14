package main

import (
	config "BannerService/internal/configs"
	"BannerService/internal/gateway"
	"BannerService/internal/server"
	"BannerService/internal/service"
	lg "BannerService/pkg/logger"
	"errors"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"net/http"
)

func main() {
	logger := lg.InitLogger()

	logger.Info.Print("Executing InitConfig.")
	if err := config.InitConfig(); err != nil {
		logger.Err.Fatalf(err.Error())
	}

	logger.Info.Print("Connecting to Postgres.")
	db, err := gateway.NewPostgresDB(viper.GetString("POSTGRES_STRING"))

	if err != nil {
		logger.Err.Fatalf(err.Error())
	}

	logger.Info.Print("Connecting to Redis.")
	redisCl, err := gateway.NewRedisClient(
		viper.GetString("REDIS_HOST"),
		viper.GetString("REDIS_PORT"),
		viper.GetString("REDIS_PASSWORD"))

	if err != nil {
		logger.Err.Fatalf(err.Error())
	}

	CacheON := viper.GetString("CACHE_ON")

	logger.Info.Print("CACHE_ON = " + CacheON)

	logger.Info.Print("Creating Gateways.")
	gateways := gateway.NewGateway(db, redisCl, CacheON == "true")

	logger.Info.Print("Creating Services.")
	services := service.NewService(gateways)

	logger.Info.Print("Creating server.")

	port := viper.GetString("PORT")
	srv := server.NewHttpServer(":"+port, services)

	logger.Info.Print("Starting the server on port: " + port + "\n\n")

	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Err.Fatalf("error occured while running http server: \"%s\" \n", err.Error())
	}
}
