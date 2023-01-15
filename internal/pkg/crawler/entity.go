package crawler

import (
	"strings"
)

type Event struct {
	Crawler  string
	Job      string
	UserName string // <firstname>.<lastname>
	UID      string // UID is determined by Crawler
	Name     string
	Message  string
}

type Channel struct {
	ID   string
	Name string
}

type Notification struct {
	Event Event
	User  Channel
}

type Post struct {
	Title  string
	Labels []string
	Bodies []Body

	Comments []Comment
}

type Comment struct {
	Bodies []Body
}

type Body struct {
	Text  string
	Files []File
}

type File struct {
	Path    string
	IsImage bool
}

// 필요하지 않을 수도 있음
func (f File) Name() string {
	splitted := strings.Split(f.Path, "/")
	return splitted[len(splitted)-1]
}
