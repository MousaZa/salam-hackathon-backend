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
	session := s.GenAi.StartChat()
	l := models.LearningRequest{
		Language:  "Python",
		Level:     "Beginner",
		FrameWork: "Django",
		Goal:      "Build a web application",
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

	fmt.Printf("Title: %s, Description: %s\n", mr.Title, mr.Description)
	for j, project := range mr.Projects {
		fmt.Printf("Projec%vt: %s, Description: %s\n", j, project.Title, project.Description)
		for i, task := range project.Tasks {
			fmt.Printf("Task%v: %s, Description: %s\n", i, task.Title, task.Description)
		}
	}

	ctx.JSON(200, gin.H{
		"data": mr,
	})
}
