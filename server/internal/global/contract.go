package global

type MessageQueue interface {
	Start() error
	Stop()
	Publish(message []byte, exchange, routingkey string, delay int64) error
	Handle()
}
