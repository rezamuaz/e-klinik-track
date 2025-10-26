package worker

import (
	"bytes"
	"context"
	"e-klinik/pkg"
	"e-klinik/pkg/constant"
	"encoding/gob"

	"time"

	"github.com/streadway/amqp"
)

// Task represents the repository used for publishing Task records.
type ProducerService struct {
	Ch *amqp.Channel
}

// NewTask instantiates the Task repository.
func NewQueueService(Ch *amqp.Channel) *ProducerService {
	return &ProducerService{
		Ch: Ch,
	}
}

// Created publishes a message indicating a task was created.
func (t *ProducerService) Create(ctx context.Context, span string, routingKey string, task any) error {
	return t.publish(span, routingKey, task)
}

// Deleted publishes a message indicating a task was deleted.
func (t *ProducerService) Deleted(ctx context.Context, id string) error {

	return t.publish("Task.Deleted", "tasks.event.deleted", id)
}

func (t *ProducerService) publish(spanName string, routingKey string, event any) error {
	var b bytes.Buffer

	// Encode event ke dalam format gob
	if err := gob.NewEncoder(&b).Encode(event); err != nil {
		return pkg.WrapError(err, pkg.ErrorCodeInternal, "failed to encode event with gob")
	}

	// Publish ke RabbitMQ
	err := t.Ch.Publish(
		constant.ExchangeName, // exchange
		routingKey,            // routing key
		true,                  // mandatory
		false,                 // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			AppId:        "tasks-rest-server",
			ContentType:  "application/x-encoding-gob",
			Body:         b.Bytes(),
			Timestamp:    time.Now(),
		},
	)
	if err != nil {
		return pkg.WrapError(err, pkg.ErrorCodeInternal, "failed to publish message to broker")
	}

	return nil

	// _, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, spanName)
	// defer span.End()

	// span.SetAttributes(
	// 	attribute.KeyValue{
	// 		Key:   semconv.MessagingSystemKey,
	// 		Value: attribute.StringValue("rabbitmq"),
	// 	},
	// 	attribute.KeyValue{
	// 		Key:   semconv.MessagingRabbitMQRoutingKeyKey,
	// 		Value: attribute.StringValue(routingKey),
	// 	},
	// )

	//-
}
