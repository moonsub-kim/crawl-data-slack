package rss

import (
	"fmt"
	"strings"

	"github.com/anaskhan96/soup"
	"github.com/mmcdole/gofeed"
)

type CrawlerOption func(*Crawler)

type Transformer interface {
	Transform(item *gofeed.Item) *gofeed.Item
	Reason() string
	String() string
}

type transformer struct {
	reason string
	parsed string
}

func (t *transformer) Reason() string {
	return t.reason
}

func (t *transformer) String() string {
	return t.parsed
}

type categoryMustContainsTransformer struct {
	transformer

	categories []string
}

func (t *categoryMustContainsTransformer) Transform(item *gofeed.Item) *gofeed.Item {
	joined := strings.Join(item.Categories, ",")
	for _, c := range t.categories {
		if strings.Contains(joined, c) {
			return nil
		}
	}

	return item
}

func WithCategoryMustContainsTransformer(categories []string) CrawlerOption {
	return func(c *Crawler) {
		c.transformers = append(c.transformers, &categoryMustContainsTransformer{categories: categories})
	}
}

type urlMustContainsTransformer struct {
	transformer

	keywords []string
}

func (t *urlMustContainsTransformer) Transform(item *gofeed.Item) *gofeed.Item {
	for _, k := range t.keywords {
		if strings.Contains(item.Link, k) {
			return nil
		}
	}

	return item
}

func WithURLMustContainsTransformer(keywords []string) CrawlerOption {
	return func(c *Crawler) {
		c.transformers = append(c.transformers, &urlMustContainsTransformer{keywords: keywords})
	}
}

type fetchRSSTransformer struct {
	transformer
	// adRegex *regexp.Regexp
}

func (t *fetchRSSTransformer) Transform(item *gofeed.Item) *gofeed.Item {
	item.Title = "" // Remove duplicated title with description
	element := soup.HTMLParse("<html>" + item.Description + "</html>")

	description := element.FullText()
	description = strings.ReplaceAll(description, "(Feed generated with FetchRSS)", "") // Remove ad text

	img := element.Find("img")
	if img.Error == nil {
		description += "\n" + img.Attrs()["src"]
	}

	item.Description = description
	return item
}

func WithFetchRSSTransformer() CrawlerOption {
	return func(c *Crawler) { c.transformers = append(c.transformers, &fetchRSSTransformer{}) }
}

type techblogPostsTransformer struct {
	transformer
}

// 회사명을 title head에 넣어준다.
func (t *techblogPostsTransformer) Transform(item *gofeed.Item) *gofeed.Item {
	item.Title = fmt.Sprintf("[%s] %s", item.Author.Name, item.Title)
	return item
}

func WithTechBlogPostsTransformer() CrawlerOption {
	return func(c *Crawler) { c.transformers = append(c.transformers, &techblogPostsTransformer{}) }
}
