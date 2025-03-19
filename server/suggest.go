package server

import (
	"context"

	"github.com/MousaZa/salam-hackathon-backend/models"
	"github.com/gin-gonic/gin"
	"github.com/google/generative-ai-go/genai"
)

func (s *Server) RequestSuggest(ctx *gin.Context) {
	sr := models.SuggestRequest{}

	if err := ctx.BindJSON(&sr); err != nil {
		ctx.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	if sr.Preference == "" || sr.Level == "" || sr.Goal == "" {
		ctx.JSON(400, gin.H{
			"error": "Missing required fields",
		})
		return
	}

	session := s.GenAi.StartChat()

	resp, err := session.SendMessage(context.Background(), genai.Text(sr.ToPrompt()))
	if err != nil {
		ctx.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(200, resp.Candidates[0].Content.Parts)
}
