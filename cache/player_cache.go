package cache

import (
	"context"
	"encoding/json"
	"psa/model"
	"time"

	"github.com/go-redis/redis/v8"
)

type PlayerCache struct {
	redisClient *redis.Client
}

func NewPlayerCache(redisClient *redis.Client) *PlayerCache {
	return &PlayerCache{
		redisClient: redisClient,
	}
}

func (c *PlayerCache) SetTopPlayers(ctx context.Context, players []models.Player) error {
	data, err := json.Marshal(players)
	if err != nil {
		return err
	}
	return c.redisClient.Set(ctx, "top_players", data, time.Minute*10).Err()
}

func (c *PlayerCache) GetTopPlayers(ctx context.Context) ([]models.Player, error) {
	result, err := c.redisClient.Get(ctx, "top_players").Result()
	if err != nil {
		return nil, err
	}

	var players []models.Player
	err = json.Unmarshal([]byte(result), &players)
	return players, err
}
