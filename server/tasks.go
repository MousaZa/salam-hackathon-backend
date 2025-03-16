package server

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
)

type Task struct {
	Id          string `json:"id"`
	Order       int    `json:"order"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
	ProjectId   string `json:"projectId"`
}

func (s *Server) GetTasks(ctx *gin.Context) {
	id := ctx.Param("id")

	r := Request{}

	if err := ctx.BindJSON(&r); err != nil {
		ctx.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	iter := s.Firestore.Client.Collection("sessions").Doc(id).Collection("tasks").Where("projectId", "==", r.Id).Documents(context.Background())
	docs, err := iter.GetAll()
	if err != nil {
		ctx.JSON(500, gin.H{
			"error": "Learning not found",
		})
		return
	}
	s.Logger.Debug("mr", docs)
	var resp []Task
	// Now you can iterate through the documents
	for _, doc := range docs {
		// Process each document
		var mr Task
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

func (s *Server) CheckTask(ctx *gin.Context) {
	id := ctx.Param("id")

	r := Request{}

	if err := ctx.BindJSON(&r); err != nil {
		ctx.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	s.Firestore.Client.Collection("sessions").Doc(id).Collection("tasks").Doc(r.Id).Update(context.Background(), []firestore.Update{firestore.Update{Path: "completed", Value: true}})
}
