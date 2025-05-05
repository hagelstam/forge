package v1

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/hagelstam/forge/internal"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var _ internal.Repository = &Repository{}

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetPosts() ([]internal.Post, error) {
	opts := options.Find().SetSort(bson.M{"date_ordered": 1})
	cursor, err := r.db.Collection("posts").Find(context.Background(), bson.M{}, opts)
	if err != nil {
		return []internal.Post{}, err
	}

	posts := []internal.Post{}
	if err = cursor.All(context.Background(), &posts); err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *Repository) CreatePost(post internal.Post) error {
	_, err := r.db.Collection("posts").InsertOne(context.Background(), post)
	return err
}

func (r *Repository) DeletePost(ID string) error {
	filter := bson.M{"_id": ID}
	opts := options.Delete().SetHint(bson.M{"_id": 1})

	result, err := r.db.Collection("posts").DeleteOne(context.Background(), filter, opts)
	if err != nil {
		return fmt.Errorf("error deleting post %s", ID.Hex())
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("post id %s doesn't exist", ID.Hex())
	}

	return nil
}

func (r *Repository) UpdatePost(post internal.Post) error {
	filter := bson.M{"_id": post.ID}
	opts := options.Update().SetUpsert(true)

	_, err := r.db.Collection("posts").UpdateOne(context.Background(), filter, bson.M{"$set": post}, opts)
	if err != nil {
		return fmt.Errorf("error updating post %s", post.ID.Hex())
	}

	return nil
}
