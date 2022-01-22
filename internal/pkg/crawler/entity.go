package crawler

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
