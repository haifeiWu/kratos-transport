package rocketmq

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tx7do/kratos-transport/broker"
)

const (
	testBroker    = "127.0.0.1:9876"
	testTopic     = "test"
	testGroupName = "CID_ONSAPI_OWNER"
)

func TestSubscribe(t *testing.T) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	ctx := context.Background()

	b := NewBroker(
		broker.OptionContext(ctx),
		WithNameServer([]string{testBroker}),
		//WithNameServerDomain(testBroker),
	)

	_ = b.Init()

	if err := b.Connect(); err != nil {
		t.Logf("cant connect to broker, skip: %v", err)
		t.Skip()
	}

	_, err := b.Subscribe(testTopic, receive,
		broker.SubscribeContext(ctx),
		broker.Queue(testGroupName),
	)
	assert.Nil(t, err)

	<-interrupt
}

func receive(_ context.Context, event broker.Event) error {
	fmt.Printf("Topic: %s Payload: %s\n", event.Topic(), string(event.Message().Body))
	//_ = event.Ack()
	return nil
}

func TestPublish(t *testing.T) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	ctx := context.Background()

	b := NewBroker(
		broker.OptionContext(ctx),
		WithNameServer([]string{testBroker}),
		//WithNameServerDomain(testBroker),
	)

	_ = b.Init()

	if err := b.Connect(); err != nil {
		t.Logf("cant connect to broker, skip: %v", err)
		t.Skip()
	}

	var msg broker.Message
	msg.Body = []byte(`{"Humidity":60, "Temperature":25}`)
	for i := 0; i < 10; i++ {
		err := b.Publish(testTopic, &msg)
		assert.Nil(t, err)
	}

	<-interrupt
}
