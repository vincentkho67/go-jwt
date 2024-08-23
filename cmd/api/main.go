package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/vincentkho67/jwt/internal/delivery/http"
	"github.com/vincentkho67/jwt/internal/middleware"
	"github.com/vincentkho67/jwt/internal/repository"
	"github.com/vincentkho67/jwt/internal/usecase"
	"github.com/vincentkho67/jwt/pkg/config"
	"github.com/vincentkho67/jwt/pkg/database"
)

func main() {
	config.LoadEnv()

	db, err := database.NewPostgresConnection()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	repo := repository.NewRepository(db)
	userUseCase := usecase.NewUserUseCase(repo)
	noteUseCase := usecase.NewNoteUseCase(repo)
	handler := http.NewHandler(userUseCase, noteUseCase)

	r := gin.Default()

	r.POST("/register", handler.Register)
	r.POST("/login", handler.Login)

	auth := r.Group("/")
	auth.Use(middleware.AuthMiddleware())
	{
		auth.POST("/notes", handler.CreateNote)
		auth.GET("/notes", handler.GetNotes)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
