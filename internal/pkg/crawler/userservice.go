package crawler

type ChannelService interface {
	GetChannels() ([]Channel, error)
}
