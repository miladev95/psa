package routes

import (
	"github.com/gin-gonic/gin"
	"psa/controller"
)

func PlayerRoutes(router *gin.Engine, playerController *controllers.PlayerController) {
	router.POST("/players", playerController.InsertOrUpdatePlayer)
	router.GET("/players/:id", playerController.GetPlayerById)
	router.GET("/players/top", playerController.GetTopPlayers)
}
