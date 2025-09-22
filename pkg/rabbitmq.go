package pkg

import (
	"e-klinik/config"
	"e-klinik/pkg/constant"
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

// RabbitMQ ...
type RabbitMQ struct {
	Conn *amqp.Connection
}

// NewRabbitMQ instantiates the RabbitMQ instances using configuration defined in environment variables.
func NewRabbit(cfg *config.Config) (*RabbitMQ, error) {
	url := fmt.Sprintf("amqp://%s:%s@%s:%s",
		cfg.RabbitMq.User,
		cfg.RabbitMq.Password,
		cfg.RabbitMq.Host,
		cfg.RabbitMq.Port)
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	if err := ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	); err != nil {
		return nil, err
	}

	return &RabbitMQ{
		Conn: conn,
	}, nil
}

// NewChannel returns a new AMQP channel with QoS settings
func (r *RabbitMQ) NewChannel() (*amqp.Channel, error) {
	ch, err := r.Conn.Channel()
	if err != nil {
		return nil, err
	}

	if err := ch.Qos(1, 0, false); err != nil {
		return nil, err
	}

	return ch, nil
}

// Close gracefully closes the connection and channel
func (r *RabbitMQ) Close() error {

	if r.Conn != nil {
		_ = r.Conn.Close()
	}

	return nil
}

func (rmq *RabbitMQ) SetupExchange(ch *amqp.Channel) error {

	err := ch.ExchangeDeclare(
		constant.ExchangeName, // exchange name
		constant.ExchangeType, // exchange type
		true,                  // durable
		false,                 // auto-deleted
		false,                 // internal
		false,                 // no-wait
		nil,                   // arguments
	)
	if err != nil {
		failOnError(err, "ch.ExchangeDeclare")
	}

	return nil
}

func (rmq *RabbitMQ) SetupQueue(ch *amqp.Channel) error {
	queue, err := ch.QueueDeclare(
		constant.QueueName, true, false, false, false, nil)
	if err != nil {
		failOnError(err, "ch.QueueDeclare")
	}

	err = ch.QueueBind(
		queue.Name,                              // queue name
		fmt.Sprintf("%s#", constant.RoutingKey), // routing key
		constant.ExchangeName,                   // exchange
		false,
		nil,
	)
	return nil
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
