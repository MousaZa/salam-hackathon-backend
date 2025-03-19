package server

import (
	"context"

	"github.com/MousaZa/salam-hackathon-backend/models"
	"github.com/gin-gonic/gin"
	"github.com/google/generative-ai-go/genai"
)

func (s *Server) RequestHelp(ctx *gin.Context) {
	sr := models.HelpRequest{}

	if err := ctx.BindJSON(&sr); err != nil {
		ctx.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	if sr.FrameWork == "" || sr.Language == "" || sr.Task == "" || sr.Project == "" {
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

	// var textContent string
	// for _, part := range resp.Candidates[0].Content.Parts {
	// 	if textPart, ok := part.(genai.Text); ok {
	// 		textContent += string(textPart)
	// 	}
	// }

	// fmt.Printf("Response: %s\n", textContent)

	// // The response is in the format: [```json {...} ```]
	// // We need to extract just the JSON part
	// jsonStart := strings.Index(textContent, "{")
	// jsonEnd := strings.LastIndex(textContent, "}")

	// if jsonStart < 0 || jsonEnd < 0 || jsonEnd <= jsonStart {
	// 	ctx.JSON(500, gin.H{
	// 		"error": "Invalid JSON response format",
	// 	})
	// 	return
	// }

	// jsonContent := textContent[jsonStart : jsonEnd+1]

	// sur := models.SuggestResponse{}
	// if err := json.Unmarshal([]byte(jsonContent), &sur); err != nil {
	// 	ctx.JSON(500, gin.H{
	// 		"error": err.Error(),
	// 	})
	// 	return
	// }

	ctx.JSON(200, resp.Candidates[0].Content.Parts)
}
