package server

import (
	"time"

	"github.com/MousaZa/salam-hackathon-backend/db"
	"github.com/gin-contrib/cors"
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

	s.Router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	s.Router.POST("/new-learning/:id", s.NewLearning)

	s.Router.POST("/projects/:id", s.GetProjects)
	s.Router.POST("/tasks/:id", s.GetTasks)
	s.Router.POST("/tasks/check/:id", s.CheckTask)
	s.Router.GET("/learnings/:id", s.GetLearnings)

	s.Router.POST("/help", s.RequestHelp)

	s.Router.POST("/suggest", s.RequestSuggest)
}
