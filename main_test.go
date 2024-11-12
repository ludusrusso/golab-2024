package main_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNatsJS(t *testing.T) {
	js := runEmbeddedServer(t)
	streamName := "orders"

	ctx := context.Background()

	createStream(t, js, streamName)
	for i := 0; i < 8; i++ {
		_, err := js.Publish(ctx, fmt.Sprintf("orders.EU.%d.shipped", i), []byte(""))
		require.NoError(t, err)
	}

	sub, err := js.CreateConsumer(ctx, streamName, jetstream.ConsumerConfig{
		Name:           "tests",
		FilterSubjects: []string{"orders.>"},
		DeliverPolicy:  jetstream.DeliverAllPolicy,
	})
	require.NoError(t, err)

	var recevedCnt int
	cctx, err := sub.Consume(func(msg jetstream.Msg) {
		recevedCnt++
		msg.Ack()
	})
	require.NoError(t, err)
	t.Cleanup(func() {
		cctx.Stop()
	})

	assert.EventuallyWithT(t, func(t *assert.CollectT) {
		assert.Equal(t, 10, recevedCnt)
	}, 5*time.Second, 100*time.Millisecond)

}

func createStream(t *testing.T, js jetstream.JetStream, streamName string) {
	t.Helper()

	ctx := context.Background()
	_, err := js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:     streamName,
		Subjects: []string{"orders.*.*.*"},
	})
	if err != nil {
		t.Fatalf("Error creating stream: %v", err)
	}
}

func runEmbeddedServer(t *testing.T) jetstream.JetStream {
	t.Helper()

	opts := &server.Options{
		ServerName:      "embedded_server",
		DontListen:      true,
		JetStream:       true,
		JetStreamDomain: "embedded",
		StoreDir:        t.TempDir(),
	}

	ns, err := server.NewServer(opts)
	if err != nil {
		t.Fatalf("Error starting server: %v", err)
	}

	go ns.Start()

	if !ns.ReadyForConnections(5 * time.Second) {
		t.Fatalf("not ready for connection: %v", err)
	}

	nc, err := nats.Connect(nats.DefaultURL, nats.InProcessServer(ns))
	if err != nil {
		t.Fatalf("Error connecting to server: %v", err)
	}

	t.Cleanup(func() {
		nc.Close()
		ns.Shutdown()
	})

	js, err := jetstream.New(nc)
	if err != nil {
		t.Fatalf("Error creating JetStream context: %v", err)
	}

	return js
}
