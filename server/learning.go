package server

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/MousaZa/salam-hackathon-backend/models"
	"github.com/gin-gonic/gin"
	"github.com/google/generative-ai-go/genai"
)

func (s *Server) NewLearning(ctx *gin.Context) {

	id := ctx.Param("id")

	session := s.GenAi.StartChat()

	l := models.LearningRequest{}

	if err := ctx.BindJSON(&l); err != nil {
		ctx.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	if l.Language == "" || l.Level == "" || l.FrameWork == "" || l.Goal == "" {
		ctx.JSON(400, gin.H{
			"error": "Missing required fields",
		})
		return
	}

	resp, err := session.SendMessage(context.Background(), genai.Text(l.ToPrompt()))
	if err != nil {
		ctx.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Extract the text content from the response
	var textContent string
	for _, part := range resp.Candidates[0].Content.Parts {
		if textPart, ok := part.(genai.Text); ok {
			textContent += string(textPart)
		}
	}

	fmt.Printf("Response: %s\n", textContent)

	// The response is in the format: [```json {...} ```]
	// We need to extract just the JSON part
	jsonStart := strings.Index(textContent, "{")
	jsonEnd := strings.LastIndex(textContent, "}")

	if jsonStart < 0 || jsonEnd < 0 || jsonEnd <= jsonStart {
		ctx.JSON(500, gin.H{
			"error": "Invalid JSON response format",
		})
		return
	}

	jsonContent := textContent[jsonStart : jsonEnd+1]

	// Now unmarshal the actual JSON content
	var mr models.LearningResponse
	if err := json.Unmarshal([]byte(jsonContent), &mr); err != nil {
		ctx.JSON(500, gin.H{
			"error": "Failed to unmarshal response: " + err.Error(),
		})
		return
	}

	if mr.Title == "" || mr.Description == "" || len(mr.Projects) == 0 {
		ctx.JSON(500, gin.H{
			"error": "Invalid response format",
		})
		return
	}

	s.Firestore.Client.Collection("learning").Doc(id).Set(context.Background(), mr.ToMap())

	fmt.Printf("Title: %s, Description: %s\n", mr.Title, mr.Description)
	for j, project := range mr.Projects {
		fmt.Printf("Projec%vt: %s, Description: %s\n", j, project.Title, project.Description)
		for i, task := range project.Tasks {
			fmt.Printf("Task%v: %s, Description: %s\n", i, task.Title, task.Description)
		}
	}

	ctx.JSON(200, mr)
}

func (s *Server) GetLearning(ctx *gin.Context) {
	id := ctx.Param("id")

	doc, err := s.Firestore.Client.Collection("learning").Doc(id).Get(context.Background())
	if err != nil {
		s.Logger.Error("Failed to get learning response", "error", err)
		ctx.JSON(500, gin.H{
			"error": "Learning not found",
		})
		return
	}

	var mr models.LearningResponse
	if err := doc.DataTo(&mr); err != nil {
		s.Logger.Error("Failed to unmarshal learning response", "error", err)
		ctx.JSON(500, gin.H{
			"error": "Failed to unmarshal learning response",
		})
		return
	}

	ctx.JSON(200, mr)
}
