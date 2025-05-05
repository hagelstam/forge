package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/hagelstam/forge/internal"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type mockRepository struct {
	internal.Repository
	getPosts   func() ([]internal.Post, error)
	createPost func(internal.Post) error
	deletePost func(primitive.ObjectID) error
	updatePost func(internal.Post) error
}

func (m mockRepository) GetPosts() ([]internal.Post, error)     { return m.getPosts() }
func (m mockRepository) CreatePost(p internal.Post) error       { return m.createPost(p) }
func (m mockRepository) DeletePost(id primitive.ObjectID) error { return m.deletePost(id) }
func (m mockRepository) UpdatePost(p internal.Post) error       { return m.updatePost(p) }

func TestGetPostsHandler(t *testing.T) {
	testObjectID := primitive.NewObjectID()

	tests := []struct {
		name             string
		mockGetPosts     func() ([]internal.Post, error)
		wantErr          bool
		expectedStatus   int
		expectedResponse []internal.Post
	}{
		{
			name: "Happy path",
			mockGetPosts: func() ([]internal.Post, error) {
				return []internal.Post{{ID: testObjectID, Content: "Test Post"}}, nil
			},
			wantErr:          false,
			expectedStatus:   http.StatusOK,
			expectedResponse: []internal.Post{{ID: testObjectID, Content: "Test Post"}},
		},
		{
			name: "Repository error",
			mockGetPosts: func() ([]internal.Post, error) {
				return nil, fmt.Errorf("database error")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := mockRepository{getPosts: tt.mockGetPosts}
			h := NewHandler(r)

			req, _ := http.NewRequest("GET", "/api/v1/posts", nil)
			rr := httptest.NewRecorder()

			h.GetPostsHandler(rr, req)

			if tt.wantErr {
				if rr.Code != http.StatusInternalServerError {
					t.Errorf("expected status 500, got %v", rr.Code)
				}
				return
			}

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %v, got %v", tt.expectedStatus, rr.Code)
			}

			var res []internal.Post
			_ = json.NewDecoder(rr.Body).Decode(&res)

			if !reflect.DeepEqual(res, tt.expectedResponse) {
				t.Errorf("expected response %v, got %v", tt.expectedResponse, res)
			}
		})
	}
}

func TestCreatePostHandler(t *testing.T) {
	tests := []struct {
		name           string
		postData       map[string]string
		mockCreatePost func(internal.Post) error
		expectedStatus int
	}{
		{
			name:           "Happy path",
			postData:       map[string]string{"content": "Test Post"},
			mockCreatePost: func(p internal.Post) error { return nil },
			expectedStatus: http.StatusCreated,
		},
		{
			name:     "Missing content",
			postData: map[string]string{"invalid_key": "Test Post"},
			mockCreatePost: func(p internal.Post) error {
				return nil
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:     "Repository error",
			postData: map[string]string{"content": "Test Post"},
			mockCreatePost: func(p internal.Post) error {
				return fmt.Errorf("database error")
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := mockRepository{createPost: tt.mockCreatePost}
			h := NewHandler(r)

			jsonData, _ := json.Marshal(tt.postData)
			req, _ := http.NewRequest("POST", "/api/v1/posts", bytes.NewBuffer(jsonData))
			rr := httptest.NewRecorder()

			h.CreatePostHandler(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %v, got %v", tt.expectedStatus, rr.Code)
			}
		})
	}
}

func TestDeletePostHandler(t *testing.T) {
	tests := []struct {
		name           string
		postID         string
		mockDeletePost func(primitive.ObjectID) error
		expectedStatus int
	}{
		{
			name:   "Happy path",
			postID: primitive.NewObjectID().Hex(),
			mockDeletePost: func(id primitive.ObjectID) error {
				return nil
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			name:   "Repository error",
			postID: primitive.NewObjectID().Hex(),
			mockDeletePost: func(id primitive.ObjectID) error {
				return fmt.Errorf("database error")
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mockRepository{deletePost: tt.mockDeletePost}
			h := NewHandler(repo)

			req, _ := http.NewRequest("DELETE", "/api/v1/posts/"+tt.postID, nil)
			rr := httptest.NewRecorder()

			router := chi.NewRouter()
			router.Delete("/api/v1/posts/{id}", h.DeletePostHandler)
			router.ServeHTTP(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %v, got %v", tt.expectedStatus, rr.Code)
			}
		})
	}
}

func TestUpdatePostHandler(t *testing.T) {
	tests := []struct {
		name           string
		postID         string
		postData       map[string]string
		mockUpdatePost func(internal.Post) error
		expectedStatus int
	}{
		{
			name:   "Happy path",
			postID: primitive.NewObjectID().Hex(),
			postData: map[string]string{
				"content": "Updated content",
			},
			mockUpdatePost: func(p internal.Post) error {
				return nil
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			name:   "Missing content",
			postID: primitive.NewObjectID().Hex(),
			postData: map[string]string{
				"invalid_key": "Updated content",
			},
			mockUpdatePost: func(p internal.Post) error {
				return nil
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "Repository error",
			postID: primitive.NewObjectID().Hex(),
			postData: map[string]string{
				"content": "Updated content",
			},
			mockUpdatePost: func(p internal.Post) error {
				return fmt.Errorf("database error")
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mockRepository{updatePost: tt.mockUpdatePost}
			h := NewHandler(repo)

			jsonData, _ := json.Marshal(tt.postData)
			req, _ := http.NewRequest("PUT", "/api/v1/posts/"+tt.postID, bytes.NewBuffer(jsonData))
			rr := httptest.NewRecorder()

			router := chi.NewRouter()
			router.Put("/api/v1/posts/{id}", h.UpdatePostHandler)
			router.ServeHTTP(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %v, got %v", tt.expectedStatus, rr.Code)
			}
		})
	}
}
