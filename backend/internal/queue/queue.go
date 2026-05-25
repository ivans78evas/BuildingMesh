package queue

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
	"construction-ar-backend/internal/logger"
	"go.uber.org/zap"
)

type Queue struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	url     string
}

func NewQueue(url string) (*Queue, error) {
	q := &Queue{url: url}
	if err := q.connect(); err != nil {
		return nil, err
	}
	return q, nil
}

func (q *Queue) connect() error {
	var err error
	q.conn, err = amqp.Dial(q.url)
	if err != nil {
		return err
	}

	q.channel, err = q.conn.Channel()
	if err != nil {
		return err
	}

	_, err = q.channel.QueueDeclare(
		"point-cloud-processing",
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	return err
}

func (q *Queue) Publish(ctx context.Context, body interface{}) error {
	data, err := json.Marshal(body)
	if err != nil {
		return err
	}

	return q.channel.PublishWithContext(ctx,
		"",                        // exchange
		"point-cloud-processing", // routing key
		false,                     // mandatory
		false,                     // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         data,
			DeliveryMode: amqp.Persistent,
		})
}

func (q *Queue) StartWorker(ctx context.Context) {
	msgs, err := q.channel.Consume(
		"point-cloud-processing",
		"",    // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		logger.Log.Error("Failed to register a consumer", zap.Error(err))
		return
	}

	go func() {
		for {
			select {
			case d := <-msgs:
				logger.Log.Info("Received a message from RabbitMQ", zap.ByteString("body", d.Body))
				// TODO: Process point cloud (decimation, optimization)
			case <-ctx.Done():
				return
			}
		}
	}()
}

func (q *Queue) Ping() error {
	if q.conn == nil || q.conn.IsClosed() {
		return fmt.Errorf("rabbitmq connection closed")
	}
	return nil
}

func (q *Queue) Close() {
	if q.channel != nil {
		q.channel.Close()
	}
	if q.conn != nil {
		q.conn.Close()
	}
}
