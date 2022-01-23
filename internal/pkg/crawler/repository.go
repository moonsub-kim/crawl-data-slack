package crawler

type Repository interface {
	SaveEvent(event Event) error
	GetChannel(userName string) (Channel, error)
	SaveChannels(users []Channel) error
}
