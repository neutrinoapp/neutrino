package messaging

type MessageProcessor interface {
	Process(mType int, m string) error
}
