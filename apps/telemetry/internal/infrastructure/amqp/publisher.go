package amqp

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	amqp091 "github.com/rabbitmq/amqp091-go"

	"telemetry/internal/domain"
)

type Publisher struct {
	conn *amqp091.Connection
	ch   *amqp091.Channel
}

type receivedPayload struct {
	RecordedAt time.Time `json:"recorded_at"`
	Unit       string    `json:"unit"`
	DeviceID   uuid.UUID `json:"device_id"`
	HouseID    uuid.UUID `json:"house_id"`
	Value      float64   `json:"value"`
}

const (
	dialAttempts = 20
	dialWait     = time.Second
)

func Open(url string) (*Publisher, error) {
	conn, err := dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		closeQuiet(conn)

		return nil, fmt.Errorf("open channel: %w", err)
	}

	_, err = ch.QueueDeclare(domain.QueueTelemetryReceived, true, false, false, false, nil)
	if err != nil {
		closeQuiet(ch)
		closeQuiet(conn)

		return nil, fmt.Errorf("declare queue: %w", err)
	}

	return &Publisher{conn: conn, ch: ch}, nil
}

func (publisher *Publisher) Close() {
	closeQuiet(publisher.ch)
	closeQuiet(publisher.conn)
}

type closer interface {
	Close() error
}

func closeQuiet(item closer) {
	if item == nil {
		return
	}

	if err := item.Close(); err != nil {
		log.Printf("amqp close: %v", err)
	}
}

func dial(url string) (*amqp091.Connection, error) {
	var last error

	for attempt := 0; attempt < dialAttempts; attempt++ {
		conn, err := amqp091.Dial(url)
		if err == nil {
			return conn, nil
		}

		last = err

		time.Sleep(dialWait)
	}

	return nil, fmt.Errorf("dial amqp: %w", last)
}

func (publisher *Publisher) PublishReceived(ctx context.Context, event domain.ReadingReceived) error {
	body, err := json.Marshal(receivedPayload{
		RecordedAt: event.RecordedAt,
		Unit:       event.Unit,
		DeviceID:   event.DeviceID,
		HouseID:    event.HouseID,
		Value:      event.Value,
	})
	if err != nil {
		log.Printf("publish %s: %v", domain.QueueTelemetryReceived, err)

		return nil
	}

	err = publisher.ch.PublishWithContext(ctx, "", domain.QueueTelemetryReceived, false, false, amqp091.Publishing{
		ContentType:  "application/json",
		DeliveryMode: amqp091.Persistent,
		Body:         body,
	})
	if err != nil {
		log.Printf("publish %s: %v", domain.QueueTelemetryReceived, err)
	}

	return nil
}
