package crawler

type Archive interface {
	CreatePost(post Post) error
	CreatePosts(post []Post) error
	CreateLabel(name string) error
	ListLabels() (map[string]struct{}, error)
}
