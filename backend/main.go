package main

import (
	"log"

	"blockchain/backend/api"
	"blockchain/backend/config"
	"blockchain/backend/model"
	"blockchain/backend/service"
)

func main() {
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Printf("Failed to load config file, using defaults: %v", err)
		cfg = &config.DefaultConfig
	}

	if err := model.InitDB(cfg.MySQL.DSN); err != nil {
		log.Fatalf("Failed to init MySQL: %v", err)
	}
	log.Println("MySQL connected")

	fabricClient, err := service.NewFabricClient(&cfg.Fabric)
	if err != nil {
		log.Fatalf("Failed to init Fabric client: %v", err)
	}
	defer fabricClient.Close()
	log.Println("Fabric Gateway connected")

	router := api.SetupRouter(cfg, fabricClient)
	log.Println("Server starting on :" + cfg.Server.Port)
	if err := router.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
