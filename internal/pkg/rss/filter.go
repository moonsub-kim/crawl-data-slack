package rss

import (
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
)

type CrawlerOption func(*Crawler)

type Filter interface {
	Filter(item *gofeed.Item) bool
	Reason() string
	String() string
}

type filter struct {
	reason string
	parsed string
}

func (f *filter) Reason() string {
	return f.reason
}

func (f *filter) String() string {
	return f.parsed
}

type categoryMustContainsFilter struct {
	filter

	categories []string
}

func (f *categoryMustContainsFilter) Filter(item *gofeed.Item) bool {
	joined := strings.Join(item.Categories, ",")
	for _, c := range f.categories {
		if strings.Contains(joined, c) {
			return false
		}
	}

	return true
}

func WithCategoryMustContainsFilter(categories []string) CrawlerOption {
	return func(c *Crawler) {
		c.filters = append(
			c.filters,
			&categoryMustContainsFilter{categories: categories},
		)
	}
}

type urlMustContainsFilter struct {
	filter

	keywords []string
}

func (f *urlMustContainsFilter) Filter(item *gofeed.Item) bool {
	for _, k := range f.keywords {
		if strings.Contains(item.Link, k) {
			return false
		}
	}

	return true
}

func WithURLMustContainsFilter(keywords []string) CrawlerOption {
	return func(c *Crawler) {
		c.filters = append(
			c.filters,
			&urlMustContainsFilter{keywords: keywords},
		)
	}
}

type recentFilter struct {
	filter

	time time.Time
}

func (f *recentFilter) Filter(item *gofeed.Item) bool {
	return item.PublishedParsed != nil && item.PublishedParsed.Before(f.time)
}

func WithRecentFilter(t time.Time) CrawlerOption {
	return func(c *Crawler) {
		c.filters = append(
			c.filters,
			&recentFilter{time: t},
		)
	}
}
