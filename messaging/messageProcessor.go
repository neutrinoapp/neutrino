package messaging

type MessageProcessor interface {
	Process(mType int, m Message) error
}
