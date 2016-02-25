package messaging

type MessageProcessor interface {
	Process(m string) error
}
