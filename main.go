package main

import (
	"os"
	"strconv"

	"github.com/qb0C80aE/clay/db"
	"github.com/qb0C80aE/clay/server"
)

func main() {

	database := db.Connect()
	s := server.Setup(database)
	port := "8080"

	if p := os.Getenv("PORT"); p != "" {
		if _, err := strconv.Atoi(p); err == nil {
			port = p
		}
	}

	s.Static("/ui", "ui")

	s.Run(":" + port)

}
