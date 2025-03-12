package server

import (
	"github.com/MousaZa/salam-hackathon-backend/db"
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/go-hclog"
)

type Server struct {
	Router    *gin.Engine
	Logger    *hclog.Logger
	Firestore *db.Firestore
}

func NewServer(r *gin.Engine, l *hclog.Logger, fs *db.Firestore) *Server {
	return &Server{
		Router:    r,
		Logger:    l,
		Firestore: fs,
	}
}

func (s *Server) Run(port string) {
	s.Router.Run(port)
}

func (s *Server) SetRoutes() {
	s.Router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})
}
