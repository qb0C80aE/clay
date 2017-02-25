package submodules

import (
	"github.com/qb0C80aE/clay/controllers"
	"github.com/qb0C80aE/clay/logics"
	"github.com/qb0C80aE/clay/models"
)

func HookSubmodules() {
	controllers.HookSubmodules()
	logics.HookSubmodules()
	models.HookSubmodules()
}
