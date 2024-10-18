package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Player struct {
	ID    primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name  string             `bson:"name" json:"name"`
	Score int                `bson:"score" json:"score"`
}
