package config

import (
	"context"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	MongoClient *mongo.Client
	RedisClient *redis.Client
	Ctx         = context.TODO()
)

func InitMongoDB() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(Ctx, clientOptions)
	if err != nil {
		panic(err)
	}

	MongoClient = client
}

func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	_, err := RedisClient.Ping(Ctx).Result()
	if err != nil {
		panic(err)
	}
}
