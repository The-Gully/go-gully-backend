package main

import (
	"os"

	"github.com/Astrasv/go-gully-backend/auth/google"
	"github.com/Astrasv/go-gully-backend/middleware"
	"github.com/Astrasv/go-gully-backend/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
	google.LoadEnvAndConnect()
}

func main() {
	google.Initialize(
		os.Getenv("GOOGLE_CLIENT_ID"),
		os.Getenv("GOOGLE_CLIENT_SECRET"),
		os.Getenv("GOOGLE_CALLBACK_URL"),
		os.Getenv("SESSION_SECRET"),
	)

	r := gin.Default()
	r.Use(middleware.CORS())

	routes.Setup(r)

	r.Run(":" + os.Getenv("PORT"))
}
