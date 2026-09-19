package main

// @AsyncAPI 3.0.0
// @Title Тёплый дом — события телеметрии
// @Version 1.0.0
// @Description telemetry публикует показание, scenario читает и решает, звать ли command (ADR-006).
// @DefaultContentType application/json
// @Server broker amqp://broker.warmhouse.example / "RabbitMQ"
func AsyncAPISpec() {}

// @Channel telemetry.received
// @ChannelDescription Новое показание после записи в pg-telemetry
// @Operation send
// @OperationID publishTelemetryReceived
// @Summary telemetry публикует показание
// @Message telemetryReceived TelemetryReceived
func PublishTelemetryReceived() {}

// @Channel telemetry.received
// @Operation receive
// @OperationID onTelemetryReceived
// @Summary scenario читает показание
// @Message telemetryReceived TelemetryReceived
func OnTelemetryReceived() {}
