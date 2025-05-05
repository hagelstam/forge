package internal

type Repository interface {
	GetPosts() ([]Post, error)
	CreatePost(post Post) error
	DeletePost(ID string) error
	UpdatePost(post Post) error
}
