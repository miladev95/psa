package main

import (
	"github.com/gin-gonic/gin"
	"psa/cache"
	"psa/config"
	"psa/controller"
	"psa/repository"
	"psa/routes"
	"psa/service"
)

func main() {
	config.InitMongoDB()
	config.InitRedis()

	playerRepo := repositories.NewPlayerRepository(config.MongoClient)
	playerCache := cache.NewPlayerCache(config.RedisClient)
	playerService := services.NewPlayerService(playerRepo, playerCache)
	playerController := controllers.NewPlayerController(playerService)

	r := gin.Default()
	routes.PlayerRoutes(r, playerController)

	r.Run()
}
