package amqp

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	amqp091 "github.com/rabbitmq/amqp091-go"

	"scenario/internal/domain"
	usecase "scenario/internal/usecases/scenario"
)

type Consumer struct {
	conn    *amqp091.Connection
	ch      *amqp091.Channel
	handler *usecase.Interactor
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

func Open(url string, handler *usecase.Interactor) (*Consumer, error) {
	conn, err := dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		closeQuiet(conn)

		return nil, fmt.Errorf("open channel: %w", err)
	}

	_, err = ch.QueueDeclare(
		domain.QueueTelemetryReceived,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		closeQuiet(ch)
		closeQuiet(conn)

		return nil, fmt.Errorf("declare queue: %w", err)
	}

	return &Consumer{conn: conn, ch: ch, handler: handler}, nil
}

func (consumer *Consumer) Close() {
	closeQuiet(consumer.ch)
	closeQuiet(consumer.conn)
}

func (consumer *Consumer) Run() error {
	msgs, err := consumer.ch.Consume(
		domain.QueueTelemetryReceived,
		"scenario",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("consume: %w", err)
	}

	for msg := range msgs {
		consumer.handle(msg)
	}

	return nil
}

func (consumer *Consumer) handle(msg amqp091.Delivery) {
	var payload receivedPayload
	if err := json.Unmarshal(msg.Body, &payload); err != nil {
		log.Printf("bad %s: %v", domain.QueueTelemetryReceived, err)
		nackQuiet(msg)

		return
	}

	event, err := domain.ParseReceived(
		payload.DeviceID,
		payload.HouseID,
		payload.Value,
		payload.Unit,
		payload.RecordedAt,
	)
	if err != nil {
		log.Printf("bad %s: %v", domain.QueueTelemetryReceived, err)
		nackQuiet(msg)

		return
	}

	decision := consumer.handler.Handle(event)
	log.Printf(
		"telemetry.received device=%s value=%.1f action=%s reason=%s",
		event.DeviceID,
		event.Value,
		decision.Action,
		decision.Reason,
	)
	ackQuiet(msg)
}

func ackQuiet(msg amqp091.Delivery) {
	if err := msg.Ack(false); err != nil {
		log.Printf("ack: %v", err)
	}
}

func nackQuiet(msg amqp091.Delivery) {
	if err := msg.Nack(false, false); err != nil {
		log.Printf("nack: %v", err)
	}
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
