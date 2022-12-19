package crawler

type Messenger interface {
	GetChannels() ([]Channel, error)
	Notify(e Notification) error
}
