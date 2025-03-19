package server

import (
	"context"

	"github.com/gin-gonic/gin"
)

type Request struct {
	Id string `json:"id"`
}

type Project struct {
	Id          string `json:"id"`
	Order       int    `json:"order"`
	Title       string `json:"title"`
	Description string `json:"description"`
	IsLocked    bool   `json:"isLocked"`
	LearningId  string `json:"learningId"`
	Progress    int    `json:"progress"`
}

func (s *Server) GetProjects(ctx *gin.Context) {
	id := ctx.Param("id")

	r := Request{}

	if err := ctx.BindJSON(&r); err != nil {
		ctx.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	iter := s.Firestore.Client.Collection("sessions").Doc(id).Collection("projects").Where("learningId", "==", r.Id).Documents(context.Background())
	docs, err := iter.GetAll()
	if err != nil {
		ctx.JSON(500, gin.H{
			"error": "Learning not found",
		})
		return
	}
	s.Logger.Debug("mr", docs)
	var resp []Project
	// Now you can iterate through the documents
	for _, doc := range docs {
		// Process each document
		var mr Project
		err := doc.DataTo(&mr)
		mr.Id = doc.Ref.ID
		if err != nil {
			s.Logger.Error("Failed to get learning response", "error", err)
			ctx.JSON(500, gin.H{
				"error": "Learning not found",
			})
		}
		iter1 := s.Firestore.Client.Collection("sessions").Doc(id).Collection("tasks").Where("projectId", "==", mr.Id).Documents(context.Background())
		docs1, err := iter1.GetAll()
		if err != nil {
			ctx.JSON(500, gin.H{
				"error": "Learning not found",
			})
			return
		}
		p := 0
		for _, doc1 := range docs1 {
			// Process each document
			var tr Task
			err := doc1.DataTo(&tr)
			if err != nil {
				s.Logger.Error("Failed to get learning response", "error", err)
				ctx.JSON(500, gin.H{
					"error": "Learning not found",
				})
			}

			if tr.Completed {
				p++
			}
		}
		mr.Progress = p

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
