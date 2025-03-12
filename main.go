package main

import (
	"github.com/MousaZa/salam-hackathon-backend/db"
	"github.com/MousaZa/salam-hackathon-backend/server"
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/go-hclog"
)

func main() {
	l := hclog.Default()
	r := gin.Default()
	fs, err := db.NewConnection()
	if err != nil {
		return
	}
	s := server.NewServer(r, &l, fs)

	s.Run(":8888")
}
