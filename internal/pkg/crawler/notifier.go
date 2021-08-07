package crawler

type Notifier interface {
	Notify(e Event) error
}
