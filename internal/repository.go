package internal

import "go.mongodb.org/mongo-driver/bson/primitive"

type Repository interface {
	GetPosts() ([]Post, error)
	CreatePost(post Post) error
	DeletePost(ID primitive.ObjectID) error
	UpdatePost(post Post) error
}
