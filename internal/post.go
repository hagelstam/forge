package internal

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Post struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Content   string             `bson:"content,omitempty" json:"content"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}
