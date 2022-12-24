package rss

import (
	"regexp"
	"strings"
	"time"

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

func (f *transformer) Reason() string {
	return f.reason
}

func (f *transformer) String() string {
	return f.parsed
}

type categoryMustContainsTransformer struct {
	transformer

	categories []string
}

func (f *categoryMustContainsTransformer) Transform(item *gofeed.Item) *gofeed.Item {
	joined := strings.Join(item.Categories, ",")
	for _, c := range f.categories {
		if strings.Contains(joined, c) {
			return nil
		}
	}

	return item
}

func WithCategoryMustContainsTransformer(categories []string) CrawlerOption {
	return func(c *Crawler) {
		c.transformers = append(
			c.transformers,
			&categoryMustContainsTransformer{categories: categories},
		)
	}
}

type urlMustContainsTransformer struct {
	transformer

	keywords []string
}

func (f *urlMustContainsTransformer) Transform(item *gofeed.Item) *gofeed.Item {
	for _, k := range f.keywords {
		if strings.Contains(item.Link, k) {
			return nil
		}
	}

	return item
}

func WithURLMustContainsTransformer(keywords []string) CrawlerOption {
	return func(c *Crawler) {
		c.transformers = append(
			c.transformers,
			&urlMustContainsTransformer{keywords: keywords},
		)
	}
}

type recentTransformer struct {
	transformer

	time time.Time
}

func (f *recentTransformer) Transform(item *gofeed.Item) *gofeed.Item {
	if item.PublishedParsed != nil && item.PublishedParsed.Before(f.time) {
		return item
	}
	return nil
}

func WithRecentTransformer(t time.Time) CrawlerOption {
	return func(c *Crawler) {
		c.transformers = append(
			c.transformers,
			&recentTransformer{time: t},
		)
	}
}

type fetchRSSTransformer struct {
	transformer

	adRegex *regexp.Regexp
}

func (f *fetchRSSTransformer) Transform(item *gofeed.Item) *gofeed.Item {
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
	return func(c *Crawler) {
		c.transformers = append(
			c.transformers,
			&fetchRSSTransformer{},
		)
	}
}
