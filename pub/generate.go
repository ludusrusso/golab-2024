package pub

import (
	"golab-2024/orders"
	"strconv"
	"sync"

	"github.com/nats-io/nats.go/jetstream"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var GenerateCMD = &cobra.Command{
	Use:   "generate",
	Short: "generate random orders",
	Run: func(cmd *cobra.Command, args []string) {
		num := 100
		if len(args) == 1 {
			n, err := strconv.Atoi(args[0])
			if err == nil {
				num = n
			}
		}
		js := mustGetNats()
		generateOrders(js, num)
	},
}

func generateOrders(js jetstream.JetStream, num int) {
	wg := sync.WaitGroup{}

	wg.Add(num)
	for i := 0; i < num; i++ {
		go func() {
			order := orders.NewOrder()
			if err := publishOrder(js, order); err != nil {
				logrus.Warnf("cannot publish order: %v", order)
			}

			wg.Done()
		}()
	}

	wg.Wait()
}
