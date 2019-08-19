package mqListeners

import "github.com/streadway/amqp"

type ListenInterface interface {
	StartListen(done <-chan  struct{})
	DoWork(msgs amqp.Delivery) error
}
