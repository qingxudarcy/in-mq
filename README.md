
# A Simple But Practical RabbitMq Client

A simple AMQP 0.9.1 client in Golang


## Demo

```golang
import mq "github.com/qingxudarcy/in-mq"
```

##### Connect Rabbitmq
```golang
dsn := mq.NewDsn(host, port, vhost, user, password)
conn := mq.NewConnetion(dsn, mq.NewExchange(exchangeName, exchangeType))
```

##### Open Up Consumers
```golang

func consumerHandler(s []byte) {
	fmt.Println(string(s))
}

queue := mq.NewQueue(exchangeName, queueName, []string{routingKey})
consumer := mq.NewConsumer(conn, queue, workers)
consumer.Start(consumerHandler)
```

##### Push Message
```golang
producer := mq.NewProducer(routingKey, []byte("in-mq"), conn)
producer.Push()   // Producers and consumers are advised not to share the same connection object.
```
## Installation

go install github.com/qingxudarcy/in-mq

## Features

V2 already supports:

- queue support the configuration of multiple routingkey
- prefetch count

v3……


## Support

For support, email darcychang@88.com or join our Slack channel.
