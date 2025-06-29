package broker

import (
	"context"
	"encoding/json"

	"github.com/nats-io/nats.go"
	"good_api_test/models"
)

type Broker interface {
	Publish(ctx context.Context, good *models.Good) error
}

type NatsBroker struct {
	conn *nats.Conn
}

func NewNatsBroker(url string) (*NatsBroker, error) {
	conn, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}
	return &NatsBroker{conn},
		nil
}

func (b *NatsBroker) Publish(ctx context.Context, good *models.Good) error {
	data, err := json.Marshal(good)
	if err != nil {
		return err
	}
	return b.conn.Publish("goods", data)
}
