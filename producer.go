package inmq

import (
	"time"

	"github.com/streadway/amqp"
)

type Producer struct {
	RoutingKey string
	Body       []byte
	Conn       *Connection
}

func NewProducer(routingKey, body string, conn *Connection) *Producer {
	return &Producer{
		RoutingKey: routingKey,
		Body:       []byte(body),
		Conn:       conn,
	}
}

// Push will push data onto the queue, and wait for a confirm.
// If no confirms are received until within the resendTimeout,
// it continuously re-sends messages until a confirm is received.
// This will block until the server sends a confirm. Errors are
// only returned if the push action itself fails, see UnsafePush.
func (producer *Producer) Push() error {
	defer producer.Conn.Close()
	select {
	case <-producer.Conn.session.ready:
	case <-time.After(waitInitDelay):
		return errNotConnected
	}
	producer.Conn.config.Logger.Printf("Push message is %s", producer.Body)
	times := 0
	err := errSendMsg
	for {
		if times >= retryTime {
			return err
		}
		times++
		err := producer.unsafePush()
		if err != nil {
			producer.Conn.config.Logger.Println("Push failed. Retrying...")
			select {
			case <-producer.Conn.session.done:
				return errAlreadyClosed
			case <-time.After(resendDelay):
			}
			continue
		}
		select {
		case confirm := <-producer.Conn.session.notifyConfirm:
			if !confirm.Ack {
				producer.Conn.config.Logger.Printf("nack message %d, body: %s", confirm.DeliveryTag, producer.Body)
				producer.Conn.config.Logger.Println("Push didn't confirm. Retrying...")
			} else {
				return nil
			}
		case <-time.After(resendDelay):
		}
	}
}

// UnsafePush will push to the queue without checking for
// confirmation. It returns an error if it fails to producer.Connect.
// No guarantees are provided for whether the server will
// recieve the message.
func (producer *Producer) unsafePush() error {
	if !producer.Conn.session.isReady {
		return errNotConnected
	}
	return producer.Conn.session.channel.Publish(
		producer.Conn.session.Exchange.Name, // Exchange
		producer.RoutingKey,                 // Routing key
		false,                               // Mandatory
		false,                               // Immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(producer.Body),
		},
	)
}
