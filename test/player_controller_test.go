package services_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"psa/model"
	"psa/service"
	"testing"
)

// Mock for PlayerRepository
type MockPlayerRepository struct {
	mock.Mock
}

func (m *MockPlayerRepository) InsertPlayer(ctx context.Context, player *models.Player) error {
	args := m.Called(ctx, player)
	return args.Error(0)
}

func (m *MockPlayerRepository) UpdatePlayer(ctx context.Context, player *models.Player) error {
	args := m.Called(ctx, player)
	return args.Error(0)
}

func (m *MockPlayerRepository) FindPlayerById(ctx context.Context, id primitive.ObjectID) (*models.Player, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.Player), args.Error(1)
}

func (m *MockPlayerRepository) GetTopPlayers(ctx context.Context, limit int) ([]models.Player, error) {
	args := m.Called(ctx, limit)
	return args.Get(0).([]models.Player), args.Error(1)
}

// Mock for PlayerCache
type MockPlayerCache struct {
	mock.Mock
}

func (m *MockPlayerCache) GetTopPlayers(ctx context.Context) ([]models.Player, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.Player), args.Error(1)
}

func (m *MockPlayerCache) SetTopPlayers(ctx context.Context, players []models.Player) error {
	args := m.Called(ctx, players)
	return args.Error(0)
}

// Test for CreatePlayer
func TestCreatePlayer(t *testing.T) {
	mockRepo := new(MockPlayerRepository)
	mockCache := new(MockPlayerCache)

	playerService := services.NewPlayerService(mockRepo, mockCache)

	player := &models.Player{
		Name:  "Player1",
		Score: 100,
	}

	mockRepo.On("InsertPlayer", mock.Anything, player).Return(nil)
	mockRepo.On("GetTopPlayers", mock.Anything, 100).Return([]models.Player{*player}, nil)
	mockCache.On("SetTopPlayers", mock.Anything, []models.Player{*player}).Return(nil)

	err := playerService.CreatePlayer(context.TODO(), player)

	assert.NoError(t, err)

	mockRepo.AssertCalled(t, "InsertPlayer", mock.Anything, player)
	mockRepo.AssertCalled(t, "GetTopPlayers", mock.Anything, 100)
	mockCache.AssertCalled(t, "SetTopPlayers", mock.Anything, []models.Player{*player})
}

func TestUpdatePlayer(t *testing.T) {
	mockRepo := new(MockPlayerRepository)
	mockCache := new(MockPlayerCache)

	playerService := services.NewPlayerService(mockRepo, mockCache)

	// Define the player to be updated
	player := &models.Player{
		ID:    primitive.NewObjectID(),
		Name:  "UpdatedPlayer",
		Score: 200,
	}

	// Mock repository behavior for UpdatePlayer
	mockRepo.On("UpdatePlayer", mock.Anything, player).Return(nil)

	// Mock cache refresh behavior
	mockRepo.On("GetTopPlayers", mock.Anything, 100).Return([]models.Player{*player}, nil)
	mockCache.On("SetTopPlayers", mock.Anything, []models.Player{*player}).Return(nil)

	// Call the UpdatePlayer method
	err := playerService.UpdatePlayer(context.TODO(), player)

	// Assert no error occurred
	assert.NoError(t, err)

	// Ensure the mock methods were called
	mockRepo.AssertCalled(t, "UpdatePlayer", mock.Anything, player)
	mockRepo.AssertCalled(t, "GetTopPlayers", mock.Anything, 100)
	mockCache.AssertCalled(t, "SetTopPlayers", mock.Anything, []models.Player{*player})
}

func TestGetTopPlayers(t *testing.T) {
	mockRepo := new(MockPlayerRepository)
	mockCache := new(MockPlayerCache)

	playerService := services.NewPlayerService(mockRepo, mockCache)

	// Simulate a cache miss by returning an empty slice and an error
	mockCache.On("GetTopPlayers", mock.Anything).Return([]models.Player{}, assert.AnError)

	// Define the top players to be returned from the repo
	players := []models.Player{
		{ID: primitive.NewObjectID(), Name: "Player1", Score: 100},
		{ID: primitive.NewObjectID(), Name: "Player2", Score: 90},
	}

	// Mock repository behavior to return top players
	mockRepo.On("GetTopPlayers", mock.Anything, 2).Return(players, nil)

	// Mock cache set behavior
	mockCache.On("SetTopPlayers", mock.Anything, players).Return(nil)

	// Call the GetTopPlayers method
	result, err := playerService.GetTopPlayers(context.TODO(), 2)

	// Assert no error occurred
	assert.NoError(t, err)

	// Assert the returned players are correct
	assert.Equal(t, players, result)

	// Ensure the mock methods were called
	mockCache.AssertCalled(t, "GetTopPlayers", mock.Anything)
	mockRepo.AssertCalled(t, "GetTopPlayers", mock.Anything, 2)
	mockCache.AssertCalled(t, "SetTopPlayers", mock.Anything, players)
}
