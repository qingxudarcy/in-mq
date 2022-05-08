package inmq

import (
	"errors"
	"time"
)

type ConsumerHandler func(d []byte)

type Logx interface {
	Printf(format string, v ...interface{})
	Println(v ...interface{})
	Panic(v ...interface{})
}

const (
	// When reconnecting to the server after connection failure
	reconnectDelay = 2 * time.Second

	// When setting up the object after a object exception
	reInitDelay = 2 * time.Second

	// When resending messages the server didn't confirm
	resendDelay = 2 * time.Second

	// The Maximum time to wait for start consumer or push message
	waitInitDelay = 5 * time.Second

	// The number of enrichment when rabbitmq object is not ready
	retryTime = 3
)

var (
	errNotConnected  = errors.New("not connected to a server")
	errAlreadyClosed = errors.New("already closed: not connected to the server")
	errShutdown      = errors.New("session is shutting down")
	errSendMsg       = errors.New("push message is fail")
)
