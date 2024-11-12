package pub

import (
	"context"
	"encoding/json"
	"golab-2024/orders"
	"os"
	"os/signal"
	"syscall"

	"math/rand"

	"github.com/nats-io/nats.go/jetstream"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var ProcessCMD = &cobra.Command{
	Use:   "process",
	Short: "process created orders",
	Run: func(cmd *cobra.Command, args []string) {
		js := mustGetNats()
		processCreatedOrders(js)
	},
}

func processCreatedOrders(js jetstream.JetStream) {
	ctx := context.Background()
	c, err := js.CreateOrUpdateConsumer(ctx, "ORDERS", jetstream.ConsumerConfig{
		Name:           "order-created-processor",
		Durable:        "order-created-processor",
		FilterSubjects: []string{"orders.*.*.created"},
		AckPolicy:      jetstream.AckExplicitPolicy,
		DeliverPolicy:  jetstream.DeliverAllPolicy,
		MaxAckPending:  1000,
		ReplayPolicy:   jetstream.ReplayInstantPolicy,
	})
	if err != nil {
		logrus.Fatalf("cannot create consumer: %v", err)
	}

	cctx, err := c.Consume(func(msg jetstream.Msg) {
		var order orders.Order
		if err := json.Unmarshal(msg.Data(), &order); err != nil {
			logrus.Errorf("cannot process msg: %v", err)
			msg.Nak()
			return
		}

		if rand.Intn(3) == 0 {
			order.Cancel()
		} else {
			order.Ship()
		}

		if err := publishOrder(js, order); err != nil {
			msg.Nak()
			return
		}

		msg.Ack()
	})
	if err != nil {
		logrus.Fatalf("cannot consume messages: %v", err)
	}

	defer cctx.Stop()

	// Gracefully shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)
	<-done
}
