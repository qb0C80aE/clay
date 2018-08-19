package server

import (
	"github.com/qb0C80aE/clay/middleware"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/qb0C80aE/clay/controller" // Install Clay controller by importing
	_ "github.com/qb0C80aE/clay/model"      // Install Clay model by importing
	"github.com/qb0C80aE/clay/router"
)

// Setup setups the server
func Setup(engine *gin.Engine, db *gorm.DB) (*gin.Engine, error) {
	engine.Use(gin.Logger())
	engine.Use(middleware.Recover())
	engine.Use(middleware.SetDBtoContext(db))
	engine.Use(middleware.PreloadBody())
	if err := router.Setup(engine); err != nil {
		return nil, err
	}
	return engine, nil
}
