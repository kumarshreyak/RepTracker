package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Email     string             `bson:"email" json:"email"`
	Name      string             `bson:"name" json:"name"`
	GoogleID  string             `bson:"google_id,omitempty" json:"google_id,omitempty"`
	Picture   string             `bson:"picture,omitempty" json:"picture,omitempty"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

type Session struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID       primitive.ObjectID `bson:"user_id" json:"user_id"`
	AccessToken  string             `bson:"access_token" json:"access_token"`
	RefreshToken string             `bson:"refresh_token,omitempty" json:"refresh_token,omitempty"`
	ExpiresAt    time.Time          `bson:"expires_at" json:"expires_at"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
}
