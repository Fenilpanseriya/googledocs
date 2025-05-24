package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Document struct {
	Id         primitive.ObjectID `bson:"_id"`
	Data       interface{}        `json:"data"`
	Created_at time.Time          `json:"created_at"`
	Updated_at time.Time          `json:"updated_at"`
}
