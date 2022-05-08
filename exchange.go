package inmq

type exchange struct {
	Name   string
	Type   string
	Option *exchangeOption
}

type exchangeOption struct {
	Durable    bool
	AutoDelete bool
	Internal   bool
	NoWait     bool
}

type ecOption struct {
	apply func(option *exchangeOption)
}

func WithEcDurable(durable bool) *ecOption {
	return &ecOption{
		apply: func(option *exchangeOption) {
			option.Durable = durable
		},
	}
}

func WithEcAutoDelete(autoDelete bool) *ecOption {
	return &ecOption{
		apply: func(option *exchangeOption) {
			option.AutoDelete = autoDelete
		},
	}
}

func WithEcInternal(internal bool) *ecOption {
	return &ecOption{
		apply: func(option *exchangeOption) {
			option.Internal = internal
		},
	}
}

func WithEcNoWait(noWait bool) *ecOption {
	return &ecOption{
		apply: func(option *exchangeOption) {
			option.NoWait = noWait
		},
	}
}

func NewExchange(name, etype string, opts ...*ecOption) *exchange {
	options := &exchangeOption{}
	for _, opt := range opts {
		opt.apply(options)
	}
	return &exchange{
		Name:   name,
		Type:   etype,
		Option: options,
	}
}
