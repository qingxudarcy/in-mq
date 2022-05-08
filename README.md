
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

consumer := mq.NewConsumer(consumerName, workers)
consumer.Start(conn, errorHandler, nil)
```

##### Push Message
```golang
producer := mq.Producer{RoutingKey: consumerName, Body: "Test Success!"}
producer.Push(conn)   // Producers and consumers are advised not to share the same connection object.
```
## Installation

go install github.com/qingxudarcy/in-mq

## Features

- queue support the configuration of multiple routingkey
- prefetch count


## Support

For support, email darcychang@88.com or join our Slack channel.
