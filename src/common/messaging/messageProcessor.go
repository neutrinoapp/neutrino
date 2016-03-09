package messaging

type MessageProcessor interface {
	Process(m string) (interface{}, error)
}
