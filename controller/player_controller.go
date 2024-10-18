package controllers

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"psa/model"
	"psa/service"
)

type PlayerController struct {
	playerService *services.PlayerService
}

func NewPlayerController(service *services.PlayerService) *PlayerController {
	return &PlayerController{
		playerService: service,
	}
}

func (pc *PlayerController) InsertOrUpdatePlayer(c *gin.Context) {
	var player models.Player

	// Bind the request JSON to the player model
	if err := c.ShouldBindJSON(&player); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if player.ID.IsZero() {
		// If no ID is provided, create a new player
		player.ID = primitive.NewObjectID() // Generate a new ObjectID
		err := pc.playerService.CreatePlayer(c.Request.Context(), &player)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, player) // Return the newly created player
	} else {
		// If an ID is provided, attempt to update the player
		err := pc.playerService.UpdatePlayer(c.Request.Context(), &player)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.JSON(http.StatusNotFound, gin.H{"error": "Player not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
			return
		}
		c.JSON(http.StatusOK, player) // Return the updated player
	}
}

func (pc *PlayerController) GetPlayerById(c *gin.Context) {
	id := c.Param("id") // Get the ID from the request URL

	// Convert string ID to MongoDB ObjectID
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	// Call the service to fetch the player by ID
	player, err := pc.playerService.GetPlayerById(c.Request.Context(), objID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Player not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, player)
}

func (pc *PlayerController) GetTopPlayers(c *gin.Context) {
	players, err := pc.playerService.GetTopPlayers(c.Request.Context(), 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, players)
}
