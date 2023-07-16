package crawler

import "time"

type Messenger interface {
	GetChannels() ([]Channel, error)
	Notify(e Notification) error
	ArchivePosts(channel Channel, fromDate time.Time, toDate time.Time) ([]Post, error)
	GetLabels() ([]string, error)
}
