package inmq

import (
	"github.com/streadway/amqp"
)

type session struct {
	isReady         bool
	done            chan bool
	ready           chan bool
	channel         *amqp.Channel
	Exchange        *exchange
	notifyConnClose chan *amqp.Error
	notifyChanClose chan *amqp.Error
	notifyConfirm   chan amqp.Confirmation
	connection      *amqp.Connection
}

func newSession() *session {
	return &session{
		notifyConnClose: make(chan *amqp.Error),
		notifyChanClose: make(chan *amqp.Error),
		notifyConfirm:   make(chan amqp.Confirmation, 1),
		done:            make(chan bool),
		ready:           make(chan bool),
	}
}

// changeConnection takes a new connection to the queue,
// and updates the close listener to reflect this.
func (session *session) changeConnection(connection *amqp.Connection) {
	session.connection = connection
	session.connection.NotifyClose(session.notifyConnClose)
}

// changeChannel takes a new channel to the queue,
// and updates the channel listeners to reflect this.
func (session *session) changeChannel(channel *amqp.Channel) {
	session.channel = channel
	session.channel.NotifyClose(session.notifyChanClose)
	session.channel.NotifyPublish(session.notifyConfirm)
}
