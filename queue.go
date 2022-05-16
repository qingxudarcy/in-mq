package inmq

import "github.com/streadway/amqp"

type queue struct {
	Exchange   string
	Name       string
	Option     *queueOption
	RoutingKey []string
	Args       amqp.Table
}

type queueOption struct {
	Durable    bool // defailt: true
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
}

var defaultQueueOption = &queueOption{
	Durable:    true,
	AutoDelete: false,
	Exclusive:  false,
	NoWait:     false,
}

type queOption struct {
	apply func(option *queueOption)
}

func WithQueDurable(durable bool) *queOption {
	return &queOption{
		apply: func(option *queueOption) {
			option.Durable = durable
		},
	}
}

func WithQueAutoDeleteOption(autoDelete bool) *queOption {
	return &queOption{
		apply: func(option *queueOption) {
			option.AutoDelete = autoDelete
		},
	}
}

func WithQueExclusive(exclusive bool) *queOption {
	return &queOption{
		apply: func(option *queueOption) {
			option.Exclusive = exclusive
		},
	}
}

func WithQueNowait(nowait bool) *queOption {
	return &queOption{
		apply: func(option *queueOption) {
			option.NoWait = nowait
		},
	}
}

func NewQueue(exchange, name string, routingKey []string, opts ...queOption) *queue {
	option := defaultQueueOption
	for _, opt := range opts {
		opt.apply(option)
	}

	args := make(amqp.Table, 1)
	args["x-queue-type"] = "classic"

	return &queue{
		Exchange:   exchange,
		Name:       name,
		RoutingKey: routingKey,
		Option:     option,
		Args:       args,
	}
}
