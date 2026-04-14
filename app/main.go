package main

import (
	"log"
	"net/http"
	"test_task/internal/api"
	"test_task/internal/repository"
	"test_task/internal/service"
)

func main() {
	repo := repository.NewInMemoryDeviceRepository()

	deviceService := service.NewDeviceService(repo)
	factory := &service.DefaultSignerFactory{}
	signService := service.NewSignService(repo, factory)

	handler := api.NewHandler(deviceService, signService)

	log.Fatal(http.ListenAndServe(":8080", handler.Routes()))
}
