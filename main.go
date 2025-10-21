package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/albkvv/student-job-finder-back/internal/application/usecases"
	"github.com/albkvv/student-job-finder-back/internal/db"
	"github.com/albkvv/student-job-finder-back/internal/infrastructure/inmemory"
	"github.com/albkvv/student-job-finder-back/internal/infrastructure/mongo"
	"github.com/albkvv/student-job-finder-back/internal/interfaces/http/handlers"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Файл .env не найден, используем системные переменные окружения")
	}

	if os.Getenv("MONGODB_URI") == "" {
		os.Setenv("MONGODB_URI", "mongodb+srv://alibekdias36_db_user:di%40s_o5@cluster0.lkhseyf.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0")
	}
	if os.Getenv("MONGODB_DB") == "" {
		os.Setenv("MONGODB_DB", "student_job_finder")
	}

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, 
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := db.GetMongoClient(ctx)
	if err != nil {
		log.Fatalf("mongo connect error: %v", err)
	}

	dbName := os.Getenv("MONGODB_DB")
	if dbName == "" {
		dbName = "app"
	}
	usersColl := client.Database(dbName).Collection("users")
	userRepo := mongo.NewMongoUserRepo(usersColl)
	codeRepo := inmemory.NewInMemoryCodeRepo()
	authService := usecases.NewAuthService(userRepo, codeRepo)
	authHandler := handlers.NewAuthHandler(authService)

	api := r.Group("/api")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status":  "ok",
				"message": "Server is running",
			})
		})

		api.POST("/request-code", authHandler.RequestCode)
		api.POST("/verify-code", authHandler.VerifyCode)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
