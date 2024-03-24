package main

import (
	"context"
	delivery "github.com/buguzei/medods-task/internal/delivery/http"
	mongo2 "github.com/buguzei/medods-task/internal/repo/mongo"
	"github.com/buguzei/medods-task/internal/usecase"
	"github.com/buguzei/medods-task/pkg/config"
	"github.com/buguzei/medods-task/server"
	"github.com/gorilla/mux"
	"log"
)

func main() {
	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	mongo := mongo2.NewMongo(ctx, cfg)

	defer func() {
		err = mongo.Client.Disconnect(ctx)
		if err != nil {
			log.Printf("error of disconnecting mongo: %v", err)
		}
	}()

	uc := usecase.NewUseCase(mongo, cfg)

	handler := delivery.NewHandler(cfg, uc)

	router := mux.NewRouter()

	s := server.HTTPServer{}

	log.Println("starting listening")
	if err = s.Run(cfg.Port, server.InitRoutes(router, handler)); err != nil {
		log.Fatal(err)
	}
}
