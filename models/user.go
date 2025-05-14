package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id            primitive.ObjectID `bson:"_id"`
	Email         string             `json:"email" validate:"email, required"`
	Password      string             `json:"password" validate:"required, min=6"`
	Refresh_token string             `json:"refresh_token"`
	Created_at    time.Time          `json:"created_at"`
	Updated_at    time.Time          `json:"updated_at"`
	Access_token  string             `json:"access_token"`
}
