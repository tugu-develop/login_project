package main

import (
	"login_project/controller"
	"login_project/infrastructure"
	"login_project/repository"
	"login_project/router"
	"login_project/usecase"
	"log"
	"context"

	"github.com/joho/godotenv"
)

func main() {
	// 環境変数読込
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// PostgreSQL接続
	db := infrastructure.NewDB() 

	// Redis接続
	redisClient := infrastructure.NewRedisClient()
	ctx := context.Background()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	// メイン処理
	userRepository := repository.NewUserRepository(db)
	sessionRepository := repository.NewSessionRepository(redisClient)
	userUsecase := usecase.NewUserUsecase(userRepository, sessionRepository)
	userController := controller.NewUserController(userUsecase)
	r := router.NewRouter(userController)
	r.Run()
}
