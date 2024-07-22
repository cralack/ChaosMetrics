package global

type MessageQueue interface {
	Start() error
	Stop()
	Publish(message []byte, delay int64) error
	Consume()
}
