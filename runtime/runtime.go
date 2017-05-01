package runtime

import (
	"os"
	"strconv"

	"fmt"
	_ "github.com/qb0C80aE/clay/buildtime" // Include program information
	"github.com/qb0C80aE/clay/db"
	"github.com/qb0C80aE/clay/extensions"
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

	os.Setenv("HOST", host)
	os.Setenv("PORT", port)

	database := db.Connect()
	s := server.Setup(database)

	s.Run(fmt.Sprintf("%s:%s", host, port))
}

func init() {
	runtime := &clayRuntime{}
	extensions.RegisterRuntime(runtime)
}
