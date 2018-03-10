package runtime

import (
	"os"
	"strconv"

	"fmt"
	_ "github.com/qb0C80aE/clay/buildtime" // Include program information
	"github.com/qb0C80aE/clay/db"
	"github.com/qb0C80aE/clay/extension"
	"github.com/qb0C80aE/clay/logging"
	"github.com/qb0C80aE/clay/server"
)

type clayRuntime struct {
}

func (clayRuntime *clayRuntime) Run() {
	host := "localhost"
	port := "8080"

	if h := os.Getenv("HOST"); h != "" {
		host = h
	}

	if p := os.Getenv("PORT"); p != "" {
		if _, err := strconv.Atoi(p); err == nil {
			port = p
		}
	}

	dbMode := os.Getenv("DB_MODE")
	database := db.Connect(dbMode)
	s := server.Setup(database)

	if err := s.Run(fmt.Sprintf("%s:%s", host, port)); err != nil {
		logging.Logger().Criticalf("failed to start: %s", err)
		os.Exit(1)
	}
}

func init() {
	runtime := &clayRuntime{}
	extension.RegisterRuntime(runtime)
}
