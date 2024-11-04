package xamqp

const (
	DelayType   = "x-delayed-message"
	DelayHeader = "x-delay" // millisecond as int/string
	Exchange    = "chao.exchange"
	Queue       = "chao.queue"
	RoutingKey  = "chao.routing"
)

type ROLE uint

const (
	Consumer ROLE = iota
	Producer
)

var RoleTag = []string{"consumer", "producer"}
