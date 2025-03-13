package main

import (
	"context"
	"log"
	"os"

	"github.com/MousaZa/salam-hackathon-backend/db"
	"github.com/MousaZa/salam-hackathon-backend/server"
	"github.com/gin-gonic/gin"
	"github.com/google/generative-ai-go/genai"
	"github.com/hashicorp/go-hclog"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

func main() {
	l := hclog.Default()
	r := gin.Default()

	err := godotenv.Load(".env")
	if err != nil {
		l.Error("Unable to get env", "error", err)
	}

	ctx := context.Background()

	apiKey := os.Getenv("GEMINI_API_KEY")

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-flash")

	model.SetTemperature(1)
	model.SetTopK(40)
	model.SetTopP(0.95)
	model.SetMaxOutputTokens(8192)
	model.ResponseMIMEType = "text/plain"

	fs, err := db.NewConnection()
	if err != nil {
		return
	}

	s := server.NewServer(r, l, fs, model)

	s.SetRoutes()

	s.Run(":8080")
}
