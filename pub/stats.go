package pub

import (
	"context"
	"encoding/json"
	"fmt"
	"golab-2024/orders"
	"time"

	"github.com/lucsky/cuid"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var StatsCMD = &cobra.Command{
	Use:   "stats",
	Short: "stats delivered orders",
	Run: func(cmd *cobra.Command, args []string) {
		loc := "EU"
		if len(args) == 1 {
			loc = args[0]
		}
		js := mustGetNats()
		stats(js, loc)
	},
}

func stats(js jetstream.JetStream, loc string) {
	ctx := context.Background()

	oneWeekAgo := time.Now().Add(-7 * 24 * time.Hour)

	c, err := js.CreateOrUpdateConsumer(ctx, "ORDERS", jetstream.ConsumerConfig{
		Name:           "order-total-counter" + cuid.New(),
		FilterSubjects: []string{fmt.Sprintf("orders.%s.*.shipped", loc)},
		AckPolicy:      jetstream.AckExplicitPolicy,
		DeliverPolicy:  jetstream.DeliverByStartTimePolicy,
		ReplayPolicy:   jetstream.ReplayInstantPolicy,
		OptStartTime:   &oneWeekAgo,
	})
	if err != nil {
		logrus.Fatalf("cannot create consumer: %v", err)
	}

	info, err := c.Info(context.Background())
	if err != nil {
		logrus.Fatalf("cannot get consumer info: %v", err)
	}

	count := info.NumPending

	var total float64 = 0

	cctx, err := c.Consume(func(msg jetstream.Msg) {
		var order orders.Order
		if err := json.Unmarshal(msg.Data(), &order); err != nil {
			logrus.Errorf("cannot process msg: %v", err)
			msg.Nak()
			return
		}

		total += order.Amount

		msg.Ack()
	}, jetstream.StopAfter(count))

	if err != nil {
		logrus.Fatalf("cannot consumer: %v", err)
	}

	<-cctx.Closed()

	logrus.Infof("%s orders count: %d", loc, count)
	logrus.Infof("Total %s orders: %.2f€", loc, total)
}
