package repositories

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"psa/model"
)

type PlayerRepository struct {
	collection *mongo.Collection
}

func NewPlayerRepository(db *mongo.Client) *PlayerRepository {
	return &PlayerRepository{
		collection: db.Database("game").Collection("players"),
	}
}

func (r *PlayerRepository) InsertPlayer(ctx context.Context, player *models.Player) error {
	_, err := r.collection.InsertOne(ctx, player)
	return err
}

func (r *PlayerRepository) UpdatePlayer(ctx context.Context, player *models.Player) error {
	result, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": player.ID},
		bson.M{"$set": player},
	)

	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments // No player found with that ID
	}

	return nil
}

func (r *PlayerRepository) FindPlayerById(ctx context.Context, id primitive.ObjectID) (*models.Player, error) {
	var player models.Player

	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&player)
	if err != nil {
		return nil, err
	}

	return &player, nil
}

func (r *PlayerRepository) GetTopPlayers(ctx context.Context, limit int) ([]models.Player, error) {
	opts := options.Find().SetSort(bson.D{{"score", -1}}).SetLimit(int64(limit))
	cursor, err := r.collection.Find(ctx, bson.D{}, opts)
	if err != nil {
		return nil, err
	}

	var players []models.Player
	err = cursor.All(ctx, &players)
	return players, err
}
