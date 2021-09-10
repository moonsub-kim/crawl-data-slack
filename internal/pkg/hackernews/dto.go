package hackernews

type DTO struct {
	ID         string `json:"id"`
	URL        string `json:"url"`
	CommentURL string `json:"comment_url"` // hackernews comments url
	Title      string `json:"title"`
	SubText    string `json:"subtext"`
}
