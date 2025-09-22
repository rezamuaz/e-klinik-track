package constant

const (
	QueueName       = "post_queue"
	ExchangeName    = "task"
	ExchangeType    = "topic"
	RoutingKey      = "tasks.event."
	PostCreated     = "post.created"
	PostUpdated     = "post.updated"
	PostDeleted     = "post.deleted"
	PostStatus      = "post.status"
	RMQConsumerName = "typesense-indexer"
)
