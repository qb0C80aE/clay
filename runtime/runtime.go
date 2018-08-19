package runtime

import (
	"os"

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
	environmentalVariableSet := extension.GetCurrentEnvironmentalVariableSet()

	host := environmentalVariableSet.GetClayHost()
	port := environmentalVariableSet.GetClayPortInt()

	db, err := dbpkg.Connect(environmentalVariableSet.GetClayDBMode())
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

	engine := gin.New()
	server, err := serverpkg.Setup(engine, db)
	if err != nil {
		logging.Logger().Criticalf("failed to start: %s", err)
		os.Exit(1)
	}

	extension.RegisterEngine(engine)

	if err := server.Run(fmt.Sprintf("%s:%d", host, port)); err != nil {
		logging.Logger().Criticalf("failed to start: %s", err)
		os.Exit(1)
	}

}

func init() {
	runtime := &clayRuntime{}
	extension.RegisterRuntime(runtime)
}
