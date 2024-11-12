package pub

import (
	"context"
	"encoding/json"
	"golab-2024/orders"
	"log"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/sirupsen/logrus"
)

func publishOrder(js jetstream.JetStream, order orders.Order) error {
	ctx := context.Background()

	orderBytes, err := json.Marshal(order)
	if err != nil {
		return err
	}

	_, err = js.Publish(ctx, order.Subject(), orderBytes)
	if err != nil {
		return err
	}

	logrus.WithField("subject", order.Subject()).Infof("📧 order published")

	return nil
}

func mustCreateStream(js jetstream.JetStream) jetstream.Stream {
	ctx := context.Background()
	s, err := js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:        "ORDERS",
		Description: "orders stream",
		Subjects:    []string{"orders.>"},
		Storage:     jetstream.FileStorage,
	})
	if err != nil {
		logrus.Fatalf("Error creating stream: %v", err)
	}

	return s
}

func mustGetNats() jetstream.JetStream {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		logrus.Fatalf("Error connecting to NATS: %v", err)
	}

	js, err := jetstream.New(nc)
	if err != nil {
		log.Fatal(err)
	}

	mustCreateStream(js)

	return js
}
