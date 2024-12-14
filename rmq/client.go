package rmq

import (
	"context"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

type Client struct {
	Url          string
	QueueName    string
	ExchangeName string
	Config       amqp.Config
}

func NewClient(Url, QueueName, ExchangeName string, Config amqp.Config) *Client {
	return &Client{Url: Url, QueueName: QueueName, ExchangeName: ExchangeName, Config: Config}
}

// необходимость ретраев и хранения до переподключения к очереди
func failOnError(err error, msg string) {
	if err != nil {
		logrus.Panicf("%s: %s", msg, err)
	}
}

func (c *Client) Send(msg string) {
	conn, err := amqp.DialConfig(c.Url, c.Config)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		c.QueueName, // name
		false,       // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)

	failOnError(err, "Failed to declare a queue")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body := msg
	err = ch.PublishWithContext(ctx,
		c.ExchangeName, // exchange
		q.Name,         // routing key
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})

	failOnError(err, "Failed to publish a message")

	logrus.Printf(" [x] Sent %s\n", body)
}
