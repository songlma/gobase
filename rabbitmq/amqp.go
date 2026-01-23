package rabbitmq

import (
	"context"
	"errors"
	"strconv"
	"sync"

	"github.com/songlma/gobase/contextz"
	"github.com/streadway/amqp"
)

type Pool struct {
	connMu sync.Mutex
	*amqp.Connection
	config Config
}

type Config struct {
	Addr     string
	Exchange string
}

func NewPool(ctx context.Context, config Config) *Pool {
	return &Pool{
		config: config,
	}
}

var ConnClosedErr = errors.New("rabbitmq conn is closed")
var PoolClosedErr = errors.New("rabbitmq pool is closed")

type Conn struct {
	*amqp.Channel
	queue    string
	rabbitmq *Pool
}

func (mq *Pool) NewConn(ctx context.Context, queue string) (*Conn, error) {
	if mq.Connection == nil {
		mq.connMu.Lock()
		conn, err := amqp.Dial(mq.config.Addr)
		if err != nil {
			mq.connMu.Unlock()
			return nil, err
		}
		mq.Connection = conn
		mq.connMu.Unlock()
	}
	channel, err := mq.Channel()
	if err != nil {
		errorLog(ctx, "rabbitmqPoolNewConn", err)
		return nil, err
	}
	return &Conn{
		Channel:  channel,
		rabbitmq: mq,
		queue:    queue,
	}, nil
}

func (mq *Pool) Close(ctx context.Context) error {
	if mq.Connection == nil || mq.IsClosed() {
		return nil
	}
	err := mq.Connection.Close()
	return err
}

func (conn *Conn) Publish(ctx context.Context, msg []byte) error {
	if conn.rabbitmq.IsClosed() {
		return PoolClosedErr
	}
	q, err := conn.QueueDeclare(
		conn.queue, // name
		true,       // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		return err
	}
	err = conn.Channel.Publish(
		conn.rabbitmq.config.Exchange, // exchange
		q.Name,                        // routing key
		false,                         // mandatory
		false,                         // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        msg,
		})

	if err != nil {
		return err
	}
	return nil
}

func (conn *Conn) Consume(ctx context.Context, Consumer string, f func(context.Context, Delivery)) error {
	if conn.rabbitmq.IsClosed() {
		return PoolClosedErr
	}
	q, err := conn.QueueDeclare(
		conn.queue, // name
		true,       // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		return err
	}
	msg, err := conn.Channel.Consume(
		q.Name,   // queue
		Consumer, // consumer
		false,    // auto-ack
		false,    // exclusive
		false,    // no-local
		false,    // no-wait
		nil,      // args
	)
	for {
		select {
		case d, ok := <-msg:
			if !ok {
				return ConnClosedErr
			}
			traceId := Consumer + "_" + strconv.Itoa(int(d.DeliveryTag))
			msgCtx, _ := contextz.SetTraceID(context.Background(), traceId)
			f(msgCtx, Delivery{d})
		}
	}
}

type Delivery struct {
	amqp.Delivery
}
