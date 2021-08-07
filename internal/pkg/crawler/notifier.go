package crawler

type Notifier interface {
	Notify(e Notification) error
}
