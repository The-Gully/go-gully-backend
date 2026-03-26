package main

import (
	"os"
	"strconv"

	"github.com/Astrasv/go-gully-backend/auth/google"
	"github.com/Astrasv/go-gully-backend/email"
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

	smtpPort, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	smtpTLS, _ := strconv.ParseBool(os.Getenv("SMTP_TLS"))
	email.Initialize(email.EmailConfig{
		Host:     os.Getenv("SMTP_HOST"),
		Port:     smtpPort,
		Username: os.Getenv("SMTP_USER"),
		Password: os.Getenv("SMTP_PASSWORD"),
		From:     os.Getenv("SMTP_FROM"),
		TLS:      smtpTLS,
	})

	r := gin.Default()
	r.Use(middleware.CORS())

	routes.Setup(r)

	r.Run(":" + os.Getenv("PORT"))
}
