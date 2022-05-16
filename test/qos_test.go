package test

import (
	"fmt"
	mq "github.com/qingxudarcy/in-mq"
	"log"
	"os"
	"testing"
	"time"
)

func qosHandler(s []byte) {
	fmt.Println(string(s))
}

func createQosConsumerConn() *mq.Connection {
	var std = log.New(os.Stderr, "", log.LstdFlags)
	dsn := mq.NewDsn(host, port, vhost, user, password)
	conn := mq.NewConnetion(dsn, mq.NewExchange(exchangeName, exchangeType), mq.WithRcHeartbeat(10*time.Second), mq.WithRcLogger(std), mq.WithRcQos(2))
	return conn
}

func createQosProducerConn() *mq.Connection {
	var std = log.New(os.Stderr, "", log.LstdFlags)
	dsn := mq.NewDsn(host, port, vhost, user, password)
	rabbitmq := mq.NewConnetion(dsn, mq.NewExchange(exchangeName, exchangeType), mq.WithRcLogger(std))
	return rabbitmq
}

func TestQosHandler(t *testing.T) {
	conn := createQosConsumerConn()
	queue := mq.NewQueue(exchangeName, queueName, []string{routingKey})
	consumer := mq.NewConsumer(conn, queue, workers, mq.WithCsAutoAck(true))
	go consumer.Start(qosHandler)
	producer := mq.NewProducer(routingKey, []byte("Test Qos Success!"), createQosProducerConn())
	for i := 0; i < 20; i++ {
		if err := producer.Push(); err != nil {
			t.Fatalf("Push message is error: %s", err)
		}
	}

	time.Sleep(5 * time.Second)
}
