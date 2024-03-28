package main

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Task struct {
	ID        primitive.ObjectID `bson:"_id"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdateAt  time.Time          `bson:"update_at"`
	Text      string             `bson:"text"`
	Complete  bool               `bson:"complete"`
}
