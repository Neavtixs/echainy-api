package main

import (
	"os"

	"github.com/Neavtixs/go-backend-template/configs"
	"github.com/Neavtixs/go-backend-template/internal"
)

func main() {
	configs.LoadEnv()
	db := configs.GetConnection()
	validate := configs.NewValidator()
	gin := configs.NewGin()
	log := configs.NewLogger()
	redis := configs.NewAccess()

	internal.Apps(&internal.AppsConfig{
		DB:       db,
		App:      gin,
		Validate: validate,
		Log:      log,
		Redis:    redis,
	})

	port := os.Getenv("BE_PORT")
	if port == "" {
		port = "8080"
	}

	gin.Run(":" + port)
}
