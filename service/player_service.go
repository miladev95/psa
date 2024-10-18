package services

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"psa/model"
)

type PlayerRepository interface {
	InsertPlayer(ctx context.Context, player *models.Player) error
	UpdatePlayer(ctx context.Context, player *models.Player) error
	FindPlayerById(ctx context.Context, id primitive.ObjectID) (*models.Player, error)
	GetTopPlayers(ctx context.Context, limit int) ([]models.Player, error)
}

type PlayerCache interface {
	GetTopPlayers(ctx context.Context) ([]models.Player, error)
	SetTopPlayers(ctx context.Context, players []models.Player) error
}

type PlayerService struct {
	playerRepo  PlayerRepository
	playerCache PlayerCache
}

func NewPlayerService(repo PlayerRepository, cache PlayerCache) *PlayerService {
	return &PlayerService{
		playerRepo:  repo,
		playerCache: cache,
	}
}

func (s *PlayerService) CreatePlayer(ctx context.Context, player *models.Player) error {
	err := s.playerRepo.InsertPlayer(ctx, player)
	if err != nil {
		return err
	}

	// After inserting, regenerate the cache with the updated player list
	return s.refreshCache(ctx)
}

func (s *PlayerService) UpdatePlayer(ctx context.Context, player *models.Player) error {
	err := s.playerRepo.UpdatePlayer(ctx, player)
	if err != nil {
		return err
	}

	// After updating, regenerate the cache with the updated player list
	return s.refreshCache(ctx)
}

func (s *PlayerService) GetPlayerById(ctx context.Context, id primitive.ObjectID) (*models.Player, error) {
	return s.playerRepo.FindPlayerById(ctx, id)
}

func (s *PlayerService) GetTopPlayers(ctx context.Context, limit int) ([]models.Player, error) {
	players, err := s.playerCache.GetTopPlayers(ctx)
	if err == nil {
		return players, nil
	}

	players, err = s.playerRepo.GetTopPlayers(ctx, limit)
	if err != nil {
		return nil, err
	}

	s.playerCache.SetTopPlayers(ctx, players)
	return players, nil
}

func (s *PlayerService) refreshCache(ctx context.Context) error {
	// Fetch the top players from the database
	players, err := s.playerRepo.GetTopPlayers(ctx, 100)
	if err != nil {
		return err
	}

	// Set the players list in the cache
	return s.playerCache.SetTopPlayers(ctx, players)
}
