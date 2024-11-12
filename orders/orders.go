package orders

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/lucsky/cuid"
	"github.com/sirupsen/logrus"
)

type OrderStatus string

const (
	StatusShipped  OrderStatus = "shipped"
	StatusCanceled OrderStatus = "canceled"
	StatusCreated  OrderStatus = "created"
)

type Location string

const (
	LocationUS Location = "US"
	LocationEU Location = "EU"
)

type Order struct {
	ID        string
	Timestamp string
	Location  Location
	Status    OrderStatus
	Amount    float64
}

func (o *Order) Ship() {
	o.Status = StatusShipped

	o.Logger().Infof("🚚 order shipped")
}

func (o *Order) Cancel() {
	o.Status = StatusCanceled

	o.Logger().Infof("❌ order canceled")
}

func (o Order) Logger() *logrus.Entry {
	return logrus.WithField("order", o.ID).WithField("amount", o.Amount)
}

func RandLocation() Location {
	if rand.Intn(2) == 0 {
		return LocationUS
	}

	return LocationEU
}

func (o Order) String() string {
	return fmt.Sprintf("Order %s/%s: %s - %.2f€", o.Location, o.ID, o.Status, o.Amount)
}

func (o Order) Subject() string {
	return fmt.Sprintf("orders.%s.%s.%s", o.Location, o.ID, o.Status)
}

func NewOrder() Order {
	o := Order{
		ID:        cuid.New(),
		Timestamp: time.Now().Format(time.RFC3339),
		Amount:    float64(rand.Intn(10000)) / 100.0,
		Status:    StatusCreated,
		Location:  RandLocation(),
	}

	o.Logger().Infof("🟢 new order created")

	return o
}
