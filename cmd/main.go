package main

import (
	"log"
	"os"

	"github.com/admalfrizi/weekly-wrapped-be/internal/config"
	"github.com/admalfrizi/weekly-wrapped-be/internal/controller"
	"github.com/admalfrizi/weekly-wrapped-be/internal/repository"
	"github.com/admalfrizi/weekly-wrapped-be/internal/router"
	"github.com/admalfrizi/weekly-wrapped-be/internal/service"
)

func main() {

	cfg := config.LoadConfig();

	db := config.InitDB(cfg.Database);
	defer db.Close()

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("CRITICAL ERROR: JWT_SECRET environment variable is not set!")
	}

	baseRepo := repository.NewBaseRepository(db)

	userRepo := repository.NewAuthRepository(baseRepo)
	authService := service.NewAuthService(userRepo, jwtSecret)
	authController := controller.NewAuthController(authService)

	r := router.SetupRouter(authController)

	log.Println("Server starting on port 8080");
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err);
	}
}