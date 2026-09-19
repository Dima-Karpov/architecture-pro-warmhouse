package main

import (
	"log"

	"github.com/gofiber/fiber/v3"

	"scenario/internal/domain"
	"scenario/internal/infrastructure/amqp"
	"scenario/internal/infrastructure/config"
	"scenario/internal/infrastructure/healthcheck"
	usecase "scenario/internal/usecases/scenario"
)

func main() {
	cfg := config.Load()
	handler := usecase.NewInteractor()

	consumer, err := amqp.Open(cfg.AMQPURL, handler)
	if err != nil {
		log.Fatalf("broker: %v", err)
	}
	defer consumer.Close()

	app := fiber.New()
	healthcheck.Register(app)

	go func() {
		if listenErr := app.Listen(cfg.Addr); listenErr != nil {
			log.Printf("listen: %v", listenErr)
		}
	}()

	log.Printf("scenario listens %s", domain.QueueTelemetryReceived)

	if err = consumer.Run(); err != nil {
		log.Printf("consume: %v", err)
	}
}
