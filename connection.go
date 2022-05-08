package inmq

import (
	"fmt"
	"net"
	"time"

	"github.com/streadway/amqp"
)

type Connection struct {
	DSN     string
	config  *ConnConfig
	session *session
}

func NewDsn(host, port, vhost, user, pass string) string {
	return fmt.Sprintf("amqp://%s:%s@%s:%s/%s", user, pass, host, port, vhost)
}

func NewConnetion(dsn string, exchange *exchange, connOpts ...*connOption) *Connection {
	connection := &Connection{
		DSN:     dsn,
		session: newSession(),
	}
	connection.session.Exchange = exchange
	connection.config = ConnCfgParse(connOpts...)
	go connection.handleReconnect()
	return connection
}

// handleReconnect will wait for a connection error on
// notifyConnClose, and then continuously attempt to reconnect.
func (connection *Connection) handleReconnect() {
	for {
		connection.session.isReady = false
		connection.config.Logger.Printf("Attempting to connect %s", connection.DSN)

		conn, err := connection.connect()

		if err != nil {
			connection.config.Logger.Println("Failed to connect. Retrying...")

			select {
			case <-connection.session.done:
				return
			case <-time.After(reconnectDelay):
			}
			continue
		}

		if done := connection.handleReInit(conn); done {
			break
		}
	}
}

// connect will create a new AMQP connection
func (connection *Connection) connect() (*amqp.Connection, error) {
	conn, err := amqp.DialConfig(connection.DSN, amqp.Config{Heartbeat: connection.config.Heartbeat, Dial: func(network, addr string) (net.Conn, error) {
		return net.DialTimeout(network, addr, connection.config.TimeOut)
	}})
	if err != nil {
		return nil, err
	}

	connection.session.changeConnection(conn)
	connection.config.Logger.Println("Connected!")
	return conn, nil
}

// handleReconnect will wait for a channel error
// and then continuously attempt to re-initialize both channels
func (connection *Connection) handleReInit(conn *amqp.Connection) bool {
	for {
		connection.session.isReady = false

		err := connection.init(conn)

		if err != nil {
			connection.config.Logger.Println("Failed to initialize channel. Retrying...")

			select {
			case <-connection.session.done:
				return true
			case <-time.After(reInitDelay):
			}
			continue
		}

		select {
		case <-connection.session.done:
			return true
		case <-connection.session.notifyConnClose:
			connection.config.Logger.Println("Connection closed. Reconnecting...")
			return false
		case <-connection.session.notifyChanClose:
			connection.config.Logger.Println("Channel closed. Re-running init...")
		}
	}
}

// init will initialize channel & declare exhange & declare queue
func (connection *Connection) init(conn *amqp.Connection) error {
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	err = ch.Confirm(false)
	if err != nil {
		connection.config.Logger.Printf("Failed to confirm channel, error: %s\n", err)
		return err
	}

	if err = ch.ExchangeDeclare(
		connection.session.Exchange.Name,              // name of the exchange
		connection.session.Exchange.Type,              // type
		connection.session.Exchange.Option.Durable,    // durable
		connection.session.Exchange.Option.AutoDelete, // delete when complete
		connection.session.Exchange.Option.Internal,   // internal
		connection.session.Exchange.Option.NoWait,     // noWait
		nil, // arguments
	); err != nil {
		connection.config.Logger.Printf("Failed to initialize exchange, error: %s\n", err)
		return err
	}

	connection.session.changeChannel(ch)
	connection.session.isReady = true
	connection.session.ready <- true
	connection.config.Logger.Println("Setup!")

	return nil
}

// Close will cleanly shutdown the channel and connection.
func (connection *Connection) Close() error {
	if !connection.session.isReady {
		return errShutdown
	}
	close(connection.session.done)
	connection.session.isReady = false
	err := connection.session.channel.Close()
	if err != nil {
		return err
	}
	err = connection.session.connection.Close()
	if err != nil {
		return err
	}
	return nil
}
