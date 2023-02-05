package naverd2

type Response struct {
	Links    []map[string]interface{} `json:"links"`
	Contents []Content                `json:"content"`
	Page     map[string]interface{}   `json:"page"`
}

type Content struct {
	Title       string `json:"postTitle"`
	Image       string `json:"postImage"`
	HTML        string `json:"postHtml"`
	PublishedAt int64  `json:"postPublishedAt"`
	Path        string `json:"url"`
}
