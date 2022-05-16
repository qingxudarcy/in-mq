package inmq

import (
	"fmt"
	"runtime/debug"
	"sync"
	"time"

	"github.com/streadway/amqp"
)

var wg sync.WaitGroup

type Consumer struct {
	Workers int
	Queue   *queue
	Option  *consumerOption
	Conn    *Connection
}

type consumerOption struct {
	AutoAck   bool
	Exclusive bool
	NoLocal   bool
	NoWait    bool
	QueueType string
}

// Type 指定队列的类型（x-queue-type）：classic、quorum(rabbitmq >= v3.8.0)、stream(rabbitmq >= v3.9.0)
var defaultconsumerOption = &consumerOption{
	AutoAck:   false,
	Exclusive: false,
	NoLocal:   false,
	NoWait:    false,
	QueueType: "classic",
}

type csOption struct {
	apply func(option *consumerOption)
}

func WithCsAutoAck(autoAck bool) *csOption {
	return &csOption{
		apply: func(option *consumerOption) {
			option.AutoAck = autoAck
		},
	}
}

func WithCsExclusive(exclusive bool) *csOption {
	return &csOption{
		apply: func(option *consumerOption) {
			option.Exclusive = exclusive
		},
	}
}

func WithCsNoLocal(noLocal bool) *csOption {
	return &csOption{
		apply: func(option *consumerOption) {
			option.NoLocal = noLocal
		},
	}
}

func WithCsNoWait(noWait bool) *csOption {
	return &csOption{
		apply: func(option *consumerOption) {
			option.NoWait = noWait
		},
	}
}

func WithCsQueueType(queueType string) *csOption {
	return &csOption{
		apply: func(option *consumerOption) {
			option.QueueType = queueType
		},
	}
}

func NewConsumer(conn *Connection, que *queue, workers int, opts ...*csOption) *Consumer {
	options := defaultconsumerOption
	for _, opt := range opts {
		opt.apply(options)
	}

	if conn.config.PrefetchCount > 0 {
		options.AutoAck = true
	}

	return &Consumer{
		Queue:   que,
		Workers: workers,
		Option:  options,
		Conn:    conn,
	}
}

// stream will continuously put queue items on the channel.
// It is required to call delivery.Ack when it has been
// suconsumeressfully processed, or delivery.Nack when it fails.
// Ignoring this will cause data to build up on the server.
func (consumer *Consumer) stream(delivery <-chan amqp.Delivery, handler ConsumerHandler) {
	defer wg.Done()
	for {
		select {
		case <-consumer.Conn.session.done:
			return
		case message := <-delivery:
			defer func() {
				err := recover()
				if err != nil {
					consumer.Conn.config.Logger.Printf("%s\n", debug.Stack())
					message.Nack(false, false)
					consumer.Conn.config.Logger.Println(err)
					wg.Add(1)
					consumer.stream(delivery, handler)
				}
			}()
			consumer.Conn.config.Logger.Printf("Receive message is %s\n", message.Body)
			handler(message.Body)
			if !consumer.Option.AutoAck {
				message.Ack(true)
			}
		}
	}
}

// StartConsumer will open specified consumers
// When queue is empty, the queue uses default configuration
func (consumer *Consumer) Start(handler ConsumerHandler) {
	consumer.Conn.config.Logger.Println("Attempting to start new consumers")

	if !consumer.Conn.session.isReady {
		select {
		case <-consumer.Conn.session.ready:
		case <-time.After(waitInitDelay):
			consumer.Conn.config.Logger.Panic(errNotConnected)
		}
	}

	que := consumer.Queue
	queName := que.Name
	routingKey := que.RoutingKey
	exchangeName := que.Exchange
	durable := que.Option.Durable
	autoDelete := que.Option.AutoDelete
	exclusive := que.Option.Exclusive
	noWait := que.Option.NoWait

	_, err := consumer.Conn.session.channel.QueueDeclare(
		queName,    // name of the queue
		durable,    // durable
		autoDelete, // delete when unused
		exclusive,  // exclusive
		noWait,     // noWait
		que.Args,   // arguments
	)
	if err != nil {
		consumer.Conn.config.Logger.Panicf("Failed to initialize queue, error: %s\n", err)
	}

	for _, key := range routingKey {
		if err := consumer.Conn.session.channel.QueueBind(
			queName,      // name of the queue
			key,          // bindingKey
			exchangeName, // sourceExchange
			noWait,       // noWait
			nil,          // arguments
		); err != nil {
			consumer.Conn.config.Logger.Panicf("Failed to bind queue, error: %s\n", err)
		}

	}

	wg.Add(consumer.Workers)
	for i := 0; i < consumer.Workers; i++ {
		delivery, err := consumer.Conn.session.channel.Consume(
			queName,                            // name
			fmt.Sprintf("%s_%d", queName, i+1), // consumerTag,
			consumer.Option.AutoAck,            // noAck
			consumer.Option.Exclusive,          // exclusive
			consumer.Option.NoLocal,            // noLocal
			consumer.Option.NoWait,             // noWait
			nil,                                // arguments
		)
		if err != nil {
			consumer.Conn.config.Logger.Panicf("Failed to initialize consume, error: %s\n", err)
		}
		go consumer.stream(delivery, handler)
	}
	wg.Wait()
}
