package mqListeners

import (
	"MQClient/logHelper"
	"github.com/streadway/amqp"
	"strconv"
	"time"
)

type QueueSecurity struct {
	UserName string
	Password string
	Host string
	Port int
}
type Queue struct {
	Name string
	Qos int
}

type BaseListener struct {
	*QueueSecurity
	*Queue
	mqConn *amqp.Connection
	mqChannel *amqp.Channel
}
func (this *BaseListener) Connect() (<-chan amqp.Delivery) {
	prefix := "mq_client"
	var connUrl = "amqp://" + this.UserName + ":" + this.Password + "@" + this.Host + ":" + strconv.Itoa(this.Port) + "/"
	var err error
	this.mqConn, err = amqp.Dial(connUrl)
	if err != nil {
		logHelper.FailOnError(prefix, err, "failed to connect mq server")
		return nil
	}
	this.mqChannel, err = this.mqConn.Channel()
	if err != nil {
		logHelper.FailOnError(prefix, err, "failed to create mq channel")
		this.mqConn.Close()
		return nil
	}
	q, err := this.mqChannel.QueueDeclare(
		this.Name + "_queue",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		this.mqChannel.Close()
		this.mqConn.Close()
		return nil
	}
	err = this.mqChannel.Qos(this.Qos, 0, true)
	if err != nil {
		this.mqChannel.Close()
		this.mqConn.Close()
		return nil
	}
	msgs, err := this.mqChannel.Consume(q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
		)
	if err != nil {
		this.mqChannel.Close()
		this.mqConn.Close()
		return nil
	}
	return msgs
}

func (this *BaseListener) ListenByConsume(doWork func(msgs amqp.Delivery) error, done <-chan struct{}) error {
	RECONNECT:
		msgs := this.Connect()
		if msgs == nil {
			time.Sleep(time.Second * 5)
			goto RECONNECT
		}
		defer this.mqConn.Close()
		defer this.mqChannel.Close()
		c := make(chan bool, 1)
		for {
			select {
			case d, ok := <- msgs:
				if !ok {
					logHelper.LogDefault("mq_client", "recive msg failed, retry...")
					time.Sleep(time.Second * 5)
					goto RECONNECT
				}
				doWork(d)
			case <-done:
				c <- false
				return nil
			}
		}
		<-c
	return nil
}
