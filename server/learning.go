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

type Learning struct {
	Id          string `json:"id"`
	Language    string `json:"language"`
	Level       string `json:"level"`
	FrameWork   string `json:"framework"`
	Goal        string `json:"goal"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

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
	mr.Language = l.Language
	mr.Level = l.Level
	mr.FrameWork = l.FrameWork
	mr.Goal = l.Goal

	lr, _, _ := s.Firestore.Client.Collection("sessions").Doc(id).Collection("learnings").Add(context.Background(), mr.ToMap())
	mr.Id = lr.ID
	for _, p := range mr.Projects {
		p.LearningId = lr.ID
		pre, _, _ := s.Firestore.Client.Collection("sessions").Doc(id).Collection("projects").Add(context.Background(), p.ToMap())
		for _, t := range p.Tasks {
			t.ProjectId = pre.ID
			s.Firestore.Client.Collection("sessions").Doc(id).Collection("tasks").Add(context.Background(), t.ToMap())
		}
	}
	ctx.JSON(200, mr)
}

func (s *Server) GetLearnings(ctx *gin.Context) {
	id := ctx.Param("id")

	iter := s.Firestore.Client.Collection("sessions").Doc(id).Collection("learnings").Documents(context.Background())
	docs, err := iter.GetAll()
	if err != nil {
		ctx.JSON(500, gin.H{
			"error": "Learning not found",
		})
		return
	}
	s.Logger.Debug("mr", docs)
	var resp []Learning
	// Now you can iterate through the documents
	for _, doc := range docs {
		// Process each document
		var mr Learning
		err := doc.DataTo(&mr)
		mr.Id = doc.Ref.ID
		if err != nil {
			s.Logger.Error("Failed to get learning response", "error", err)
			ctx.JSON(500, gin.H{
				"error": "Learning not found",
			})
		}
		s.Logger.Debug("mr", mr)
		resp = append(resp, mr)
		s.Logger.Debug("resp", resp)
		// ...
	}
	// if err != nil {
	// 	s.Logger.Error("Failed to get learning response", "error", err)
	// 	ctx.JSON(500, gin.H{
	// 		"error": "Learning not found",
	// 	})
	// 	return
	// }

	// if err := doc.DataTo(&mr); err != nil {
	// 	s.Logger.Error("Failed to unmarshal learning response", "error", err)
	// 	ctx.JSON(500, gin.H{
	// 		"error": "Failed to unmarshal learning response",
	// 	})
	// 	return
	// }

	ctx.JSON(200, resp)
}
