package test

var (
	exchangeName string = "test_exchange"
	exchangeType string = "direct"

	host     string = "localhost"
	port     string = "5672"
	vhost    string = "/in-mq"
	user     string = ""
	password string = ""

	queueName  string = "test_queue"
	routingKey string = "in-test"
	workers    int    = 2
)
