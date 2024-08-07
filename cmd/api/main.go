package main

import (
	"fmt"
	"github.com/antibomberman/mego-api/internal/adapters/api"
	"github.com/antibomberman/mego-api/internal/clients"
	"github.com/antibomberman/mego-api/internal/config"
	"log"
)

func main() {
	cfg := config.Load()

	authClient, err := clients.NewAuthClient(cfg.AuthServiceAddress)
	if err != nil {
		log.Printf("failed to create auth client: %v", err)
	}
	userClient, err := clients.NewUserClient(cfg.UserServiceAddress)
	if err != nil {
		log.Printf("failed to create user client: %v", err)
	}
	postClient, err := clients.NewPostClient(cfg.PostServiceAddress)
	if err != nil {
		log.Printf("failed to create post client: %v", err)
	}
	storageClient, err := clients.NewStorageClient(cfg.StorageServiceAddress)
	if err != nil {
		log.Printf("failed to create storage client: %v", err)
	}

	fmt.Println("Starting API server on port", cfg.ApiServiceServerPort)
	httpSrv := api.NewServer(cfg, authClient, userClient, postClient, storageClient)
	err = httpSrv.Start(cfg.ApiServiceServerPort)
	if err != nil {
		log.Fatalf("failed to start API server: %v", err)
	}

}
