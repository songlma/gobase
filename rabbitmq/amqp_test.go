package rabbitmq

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/songlma/gobase/logger"
)

func TestNewMQ(t *testing.T) {
	ctx := context.Background()
	mq := NewPool(ctx, Config{
		Addr: "amqp://api:api_user@127.0.0.1:5672/api",
	})
	defer func() {
		err := mq.Close(ctx)
		if err != nil {
			t.Error(err)
		}
	}()
}

func TestRabbitmq_NewConn(t *testing.T) {
	ctx := context.Background()
	mq := NewPool(ctx, Config{
		Addr: "amqp://api:api_user@127.0.0.1:5672/api",
	})
	defer func() {
		err := mq.Close(ctx)
		if err != nil {
			t.Error(err)
		}
	}()
	newConn, err := mq.NewConn(ctx, "test_q_multi_go")
	if err != nil {
		t.Error(err)
		return
	}
	err = newConn.Close()
	if err != nil {
		t.Error(err)
		return
	}
}

func TestConn_Publish(t *testing.T) {
	ctx := context.Background()
	mq := NewPool(ctx, Config{
		Addr: "amqp://api:api_user@127.0.0.1:5672/api",
	})
	defer func() {
		err := mq.Close(ctx)
		if err != nil {
			t.Error(err)
		}
	}()
	newConn, err := mq.NewConn(ctx, "test_q_multi_go")
	if err != nil {
		t.Error(err)
		return
	}
	defer func() {
		err = newConn.Close()
		if err != nil {
			t.Error(err)
			return
		}
	}()
	err = newConn.Publish(ctx, []byte("测试小时"))
	if err != nil {
		t.Error(err)
		return
	}
}

func TestConn_Consume(t *testing.T) {
	ctx := context.Background()
	mq := NewPool(ctx, Config{
		Addr: "amqp://api:api_user@127.0.0.1:5672/api",
	})
	defer func() {
		err := mq.Close(ctx)
		if err != nil {
			t.Error(err)
		}
	}()
	newConn, err := mq.NewConn(ctx, "test_q_multi_go")
	if err != nil {
		t.Error(err)
		return
	}
	defer func() {
		err = newConn.Close()
		if err != nil {
			t.Error(err)
			return
		}
	}()
	err = newConn.Consume(ctx, "test_q_multi_go", func(ctx context.Context, delivery Delivery) {
		logger.Info(ctx, string(delivery.Body))
		delivery.Ack(false)
	})
	if err != nil {
		t.Error(err)
		return
	}
}

func TestConsume_Rabbitmq_Close(t *testing.T) {
	ctx := context.Background()
	mq := NewPool(ctx, Config{
		Addr: "amqp://api:api_user@127.0.0.1:5672/api",
	})
	defer func() {
		err := mq.Close(ctx)
		if err != nil {
			t.Error(err)
		}
	}()
	newConn, err := mq.NewConn(ctx, "test_q_multi_go")
	if err != nil {
		t.Error(err)
		return
	}
	defer func() {
		err = newConn.Close()
		if err != nil {
			t.Error(err)
			return
		}
	}()
	go func() {
		select {
		case <-time.After(5 * time.Second):
			mq.Close(ctx)
		}
	}()
	err = newConn.Consume(ctx, "test_q_multi_go", func(ctx context.Context, delivery Delivery) {
		logger.Info(ctx, string(delivery.Body))
		delivery.Ack(false)
	})
	logger.Info(ctx, errors.Is(err, ConnClosedErr))
	if err != nil {
		t.Error(err)
		return
	}

}
