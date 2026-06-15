package broker

import (
	"context"
	"encoding/json"

	"github.com/nats-io/nats.go"
	"github.com/pkg/errors"
)

// Subjects for cross-service events.
const (
	SubjectOrderCreated     = "order.created"
	SubjectOrderConfirmed   = "order.confirmed"
	SubjectOrderCancelled   = "order.cancelled"
	SubjectPaymentSucceeded = "payment.succeeded"
	SubjectPaymentFailed    = "payment.failed"
)

type Broker struct {
	conn *nats.Conn
}

func New(url string) (*Broker, error) {
	conn, err := nats.Connect(url)
	if err != nil {
		return nil, errors.Wrap(err, "broker.New")
	}
	return &Broker{conn: conn}, nil
}

func (b *Broker) Publish(_ context.Context, subject string, payload interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return errors.Wrap(err, "broker.Publish.Marshal")
	}
	return errors.Wrap(b.conn.Publish(subject, data), "broker.Publish")
}

func (b *Broker) Subscribe(subject string, handler func([]byte)) error {
	_, err := b.conn.Subscribe(subject, func(msg *nats.Msg) {
		handler(msg.Data)
	})
	return errors.Wrap(err, "broker.Subscribe")
}

func (b *Broker) Close() {
	b.conn.Drain()
}
