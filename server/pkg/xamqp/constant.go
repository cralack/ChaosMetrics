package xamqp

const (
	DelayType   = "x-delayed-message"
	DelayHeader = "x-delay" // millisecond as int/string

	// Worker configuration
	WorkerRole     = "worker"
	DataExchange   = "data_exchange"
	DataQueue      = "data_queue"
	DataRoutingKey = "data_routing_key"

	// Master configuration
	MasterRole     = "master"
	TaskExchange   = "task_exchange"
	TaskQueue      = "task_queue"
	TaskRoutingKey = "task_routing_key"

	// Router configuration
	RouterRole        = "router"
	RequestExchange   = "request_exchange"
	RequestQueue      = "request_queue"
	RequestRoutingKey = "request_routing_key"

	// Analyzer configuration
	AnalyzerRole      = "analyzer"
	AnalyzeExchange   = "analyze_exchange"
	AnalyzeQueue      = "analyze_queue"
	AnalyzeRoutingKey = "analyze_routing_key"

	// Updater configuration
	UpdaterRole      = "updater"
	UpdateExchange   = "update_exchange"
	UpdateQueue      = "update_queue"
	UpdateRoutingKey = "update_routing_key"
)

type ROLE uint

const (
	Worker ROLE = iota
	Master
	Router
	Analyzer
	Updater
)

type DOM uint

const (
	Tag DOM = iota
	Exchange
	Queue
	RoutingKey
)

var Constant = [][]string{
	// Worker
	{WorkerRole, DataExchange, DataQueue, DataRoutingKey},
	// Master
	{MasterRole, TaskExchange, TaskQueue, TaskRoutingKey},
	// Router
	{RouterRole, RequestExchange, RequestQueue, RequestRoutingKey},
	// Analyzer
	{AnalyzerRole, AnalyzeExchange, AnalyzeQueue, AnalyzeRoutingKey},
	// Updater
	{UpdaterRole, UpdateExchange, UpdateQueue, UpdateRoutingKey},
}
