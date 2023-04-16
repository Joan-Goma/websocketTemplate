package websocketTemplate

import (
	"github.com/Joan-Goma/websocketTemplate/controller"
	engine "github.com/JoanGTSQ/api"
	"github.com/gin-gonic/gin"
)

const version = "DRAGON DEV 0.0.1"

func InitConfigs(DB controller.DbConnection, debugDragon, debugDB bool, router *gin.Engine, certRoute, keyRoute, port string) {

	gin.SetMode(gin.ReleaseMode)

	go controller.ReadInput(debugDragon)
	controller.InitDatabase(DB, debugDB)
	if port == ":443" {
		engine.Info.Println("Running SSL server on port :443")
		if err := router.RunTLS(":443", certRoute, keyRoute); err != nil {
			engine.Error.Fatalln("Error starting web server", err)
		}
	}

	engine.Info.Println("Running non SSL server on port", port)
	if err := router.Run(port); err != nil {
		engine.Error.Fatalln("Error starting web server", err)
	}

}
