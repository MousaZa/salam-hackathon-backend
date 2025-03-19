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

	// Get the task to check
	docRef := s.Firestore.Client.Collection("sessions").Doc(id).Collection("tasks").Doc(r.Id)
	doc, err := docRef.Get(context.Background())
	if err != nil {
		ctx.JSON(500, gin.H{
			"error": "Task not found",
		})
		return
	}

	// Convert to Task struct
	var task Task
	if err := doc.DataTo(&task); err != nil {
		ctx.JSON(500, gin.H{
			"error": "Failed to unmarshal task",
		})
		return
	}

	// Mark task as completed
	_, err = docRef.Update(context.Background(), []firestore.Update{
		{Path: "completed", Value: true},
	})
	if err != nil {
		ctx.JSON(500, gin.H{
			"error": "Failed to update task",
		})
		return
	}

	// Check if all tasks in the project are completed
	tasksRef := s.Firestore.Client.Collection("sessions").Doc(id).Collection("tasks").Where("projectId", "==", task.ProjectId)
	taskDocs, err := tasksRef.Documents(context.Background()).GetAll()
	if err != nil {
		s.Logger.Error("Failed to get project tasks", "error", err)
		ctx.JSON(500, gin.H{
			"error": "Failed to get project tasks",
		})
		return
	}

	allCompleted := true
	for _, taskDoc := range taskDocs {
		var t Task
		if err := taskDoc.DataTo(&t); err != nil {
			continue
		}
		if !t.Completed {
			allCompleted = false
			break
		}
	}

	// If all tasks are completed, unlock the next project
	if allCompleted {
		s.Logger.Info("All tasks completed in project", "projectId", task.ProjectId)

		// Get the current project to find its order
		projectRef := s.Firestore.Client.Collection("sessions").Doc(id).Collection("projects").Doc(task.ProjectId)
		projectDoc, err := projectRef.Get(context.Background())
		if err != nil {
			s.Logger.Error("Failed to get project", "error", err)
			ctx.JSON(200, gin.H{
				"message": "Task completed, but failed to get project",
			})
			return
		}

		var currentProject Project
		if err := projectDoc.DataTo(&currentProject); err != nil {
			s.Logger.Error("Failed to unmarshal project", "error", err)
			ctx.JSON(200, gin.H{
				"message": "Task completed, but failed to unmarshal project",
			})
			return
		}

		// Find the next project by order
		nextProjectsRef := s.Firestore.Client.Collection("sessions").Doc(id).Collection("projects").Where("order", "==", currentProject.Order+1)
		nextProjectDocs, err := nextProjectsRef.Documents(context.Background()).GetAll()
		if err != nil || len(nextProjectDocs) == 0 {
			// No next project or error
			ctx.JSON(200, gin.H{
				"message": "Task completed, no next project to unlock",
			})
			return
		}

		// Unlock the next project
		nextProjectRef := nextProjectDocs[0].Ref
		_, err = nextProjectRef.Update(context.Background(), []firestore.Update{
			{Path: "isLocked", Value: false},
		})
		if err != nil {
			s.Logger.Error("Failed to unlock next project", "error", err)
			ctx.JSON(200, gin.H{
				"message": "Task completed, but failed to unlock next project",
			})
			return
		}

		ctx.JSON(200, gin.H{
			"message": "Task completed and next project unlocked",
		})
		return
	}

	ctx.JSON(200, gin.H{
		"message": "Task completed",
	})
}
