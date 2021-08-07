package crawler

type Crawler interface {
	Crawl() ([]Event, error)
}
