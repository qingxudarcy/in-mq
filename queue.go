package inmq

type queue struct {
	Exchange   string
	Name       string
	RoutingKey string
	Option     *queueOption
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

func NewQueue(exchange, name, routingKey string, opts ...queOption) *queue {
	option := defaultQueueOption
	for _, opt := range opts {
		opt.apply(option)
	}
	return &queue{
		Exchange:   exchange,
		Name:       name,
		RoutingKey: routingKey,
		Option:     option,
	}
}
