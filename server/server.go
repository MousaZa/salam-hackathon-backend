package server

import (
	"github.com/MousaZa/salam-hackathon-backend/db"
	"github.com/gin-gonic/gin"
	"github.com/google/generative-ai-go/genai"
	"github.com/hashicorp/go-hclog"
)

type Server struct {
	Router    *gin.Engine
	Logger    hclog.Logger
	Firestore *db.Firestore
	GenAi     *genai.GenerativeModel
}

func NewServer(r *gin.Engine, l hclog.Logger, fs *db.Firestore, gai *genai.GenerativeModel) *Server {
	return &Server{
		Router:    r,
		Logger:    l,
		Firestore: fs,
		GenAi:     gai,
	}
}

func (s *Server) Run(port string) {
	s.Router.Run(port)
}

func (s *Server) SetRoutes() {
	s.Router.GET("/test", s.NewLearning)
}
