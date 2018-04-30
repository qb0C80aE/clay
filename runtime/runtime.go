package runtime

import (
	"os"
	"strconv"

	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/qb0C80aE/clay/buildtime" // Include program information
	dbpkg "github.com/qb0C80aE/clay/db"
	"github.com/qb0C80aE/clay/extension"
	"github.com/qb0C80aE/clay/logging"
	serverpkg "github.com/qb0C80aE/clay/server"
)

type clayRuntime struct {
}

func (clayRuntime *clayRuntime) Run() {
	host := "localhost"
	port := "8080"

	if h := os.Getenv("CLAY_HOST"); h != "" {
		host = h
	}

	if p := os.Getenv("CLAY_PORT"); p != "" {
		if _, err := strconv.Atoi(p); err == nil {
			port = p
		}
	}

	dbMode := os.Getenv("CLAY_DB_MODE")
	db, err := dbpkg.Connect(dbMode)
	if err != nil {
		logging.Logger().Critical(err.Error())
		os.Exit(1)
	}

	// Caution: Even if you input the inconsistent data like foreign keys do not exist,
	//          it will be registered, and never be checked this time.
	//          Todo: It requires order resolution logic like "depends on" between models.
	initializerList := extension.GetRegisteredInitializerList()
	modelList := extension.GetRegisteredModelList()

	db, err = extension.SetupModel(db, initializerList, modelList)
	if err != nil {
		logging.Logger().Critical(err.Error())
		os.Exit(1)
	}

	engine := gin.Default()
	server, err := serverpkg.Setup(engine, db)
	if err != nil {
		logging.Logger().Criticalf("failed to start: %s", err)
		os.Exit(1)
	}

	if err := server.Run(fmt.Sprintf("%s:%s", host, port)); err != nil {
		logging.Logger().Criticalf("failed to start: %s", err)
		os.Exit(1)
	}
}

func init() {
	runtime := &clayRuntime{}
	extension.RegisterRuntime(runtime)
}
