package test

import (
	"fmt"
	mq "github.com/qingxudarcy/in-mq"
	"log"
	"os"
	"testing"
	"time"
)

func errorHandler(s []byte) {
	fmt.Println(string(s))
	var a []int
	a[2] = 1
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

func TestErrorHandler(t *testing.T) {
	conn := createConsumerConn()
	consumer := mq.NewConsumer(conn, consumerName, workers)
	go consumer.Start(errorHandler, nil)
	producer := mq.NewProducer(consumerName, "Test Success!", createProducerConn())
	if err := producer.Push(); err != nil {
		t.Fatalf("Push message is error: %s", err)
	}
	time.Sleep(5 * time.Second)
}

func TestCloseConsumer(t *testing.T) {
	conn := createConsumerConn()
	consumer := mq.NewConsumer(conn, consumerName, workers)
	go consumer.Start(errorHandler, nil)
	time.Sleep(10 * time.Second)
	if err := conn.Close(); err != nil {
		t.Fatalf("close consumer is err %s", err)
	}
}
