package queue

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitClient struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func ConnectRabbitMQ(username, password, host, vhost string) (*amqp.Connection, error) {
	return amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s/%s", username, password, host, vhost))
}

func NewRabbitClient(conn *amqp.Connection) (*RabbitClient, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &RabbitClient{
		conn: conn,
		ch:   ch,
	}, nil
}

func (c *RabbitClient) Close() error {
	return c.ch.Close()
}

func (c *RabbitClient) CreateQueue(queueName string, durable, autodelete bool) error {
	_, err := c.ch.QueueDeclare(
		queueName,
		durable,
		autodelete,
		false,
		false,
		nil,
	)
	return err
}

func (c *RabbitClient) CreateExchange(name, kind string, durable, autodelete bool) error {
	return c.ch.ExchangeDeclare(
		name,
		kind,
		durable,
		autodelete,
		false,
		false,
		nil,
	)
}

func (c *RabbitClient) CreateBinding(queueName, exchangeName, routingKey string) error {
	return c.ch.QueueBind(
		queueName,
		routingKey,
		exchangeName,
		false,
		nil,
	)
}

// Publish payloads onto the exchange using the given routing key
func (c *RabbitClient) Send(ctx context.Context, exchangeName, routingKey string, options amqp.Publishing) error {
	return c.ch.PublishWithContext(ctx,
		exchangeName,
		routingKey,
		// Mandatory: used to determine if an error should be returned upon failure
		true,
		false,
		options)
}

func (c *RabbitClient) Consume(queue, consumer string, autoAck bool) (<-chan amqp.Delivery, error) {
	return c.ch.Consume(
		queue,
		consumer,
		autoAck,
		false,
		false,
		false,
		nil,
	)
}
