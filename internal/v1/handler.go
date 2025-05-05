package v1

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/hagelstam/forge/internal"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Handler struct {
	repository internal.Repository
}

func NewHandler(repository internal.Repository) *Handler {
	return &Handler{repository: repository}
}

func (h *Handler) GetPostsHandler(w http.ResponseWriter, r *http.Request) {
	posts, err := h.repository.GetPosts()
	if err != nil {
		http.Error(w, "error fetching posts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(posts)
}

func (h *Handler) CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	var post internal.Post

	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusInternalServerError)
		return
	}

	if post.Content == "" {
		http.Error(w, "content field is required", http.StatusBadRequest)
		return
	}

	if err := h.repository.CreatePost(post); err != nil {
		http.Error(w, "error creating post", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	postID := chi.URLParam(r, "id")

	postObjectID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		http.Error(w, "invalid post ID", http.StatusBadRequest)
		return
	}

	if err := h.repository.DeletePost(postObjectID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) UpdatePostHandler(w http.ResponseWriter, r *http.Request) {
	postID := chi.URLParam(r, "id")

	postObjectID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		http.Error(w, "invalid post ID", http.StatusBadRequest)
		return
	}

	var post internal.Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, "invalid request body", http.StatusInternalServerError)
		return
	}
	post.ID = postObjectID

	if post.Content == "" {
		http.Error(w, "content field is required", http.StatusBadRequest)
		return
	}

	if err := h.repository.UpdatePost(post); err != nil {
		http.Error(w, "error updating post", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
