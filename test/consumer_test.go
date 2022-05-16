package test

import (
	"fmt"
	mq "github.com/qingxudarcy/in-mq"
	"log"
	"os"
	"testing"
	"time"
)

func printHandler(s []byte) {
	fmt.Println(string(s))
}

func createConsumerConn() *mq.Connection {
	var std = log.New(os.Stderr, "", log.LstdFlags)
	dsn := mq.NewDsn(host, port, vhost, user, password)
	conn := mq.NewConnetion(dsn, mq.NewExchange(exchangeName, exchangeType), mq.WithRcHeartbeat(10*time.Second), mq.WithRcLogger(std))
	return conn
}

func createProducerConn() *mq.Connection {
	var std = log.New(os.Stderr, "", log.LstdFlags)
	dsn := mq.NewDsn(host, port, vhost, user, password)
	rabbitmq := mq.NewConnetion(dsn, mq.NewExchange(exchangeName, exchangeType), mq.WithRcLogger(std))
	return rabbitmq
}

func TestPrintHandler(t *testing.T) {
	conn := createConsumerConn()
	queue := mq.NewQueue(exchangeName, queueName, []string{routingKey})
	consumer := mq.NewConsumer(conn, queue, workers)
	go consumer.Start(printHandler)
	producer := mq.NewProducer(routingKey, []byte("Test Success!"), createProducerConn())
	if err := producer.Push(); err != nil {
		t.Fatalf("Push message is error: %s", err)
	}
	time.Sleep(5 * time.Second)
}
