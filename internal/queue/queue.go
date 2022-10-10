package queue

import (
	"github.com/olahol/melody"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"github.com/wagslane/go-rabbitmq"

	"github.com/carfloresf/financial-chat/config"
)

type Client struct {
	Publisher *rabbitmq.Publisher
	Consumer  *rabbitmq.Consumer
	Websocket *melody.Melody
	sessions  map[string]*melody.Session
}

type Queuer interface {
	Publish(body []byte, correlationID, exchange string) error
	Consume(consumerName, queueName, exchange string, action func(delivery rabbitmq.Delivery) rabbitmq.Action) error
	StoreSession(key string, session *melody.Session)
	GetSession(key string) *melody.Session
	RemoveSession(key string)
}

func NewClient(config *config.Config) (*Client, error) {
	publisher, err := rabbitmq.NewPublisher(
		config.RabbitMQ.URL,
		rabbitmq.Config{},
		rabbitmq.WithPublisherOptionsLogging,
	)
	if err != nil {
		log.Errorf("error creating Publisher: %s", err)

		return nil, err
	}

	consumer, err := rabbitmq.NewConsumer(
		config.RabbitMQ.URL,
		rabbitmq.Config{},
		rabbitmq.WithConsumerOptionsLogging,
	)
	if err != nil {
		log.Errorf("error creating Consumer: %s", err)

		return nil, err
	}

	return &Client{
		Publisher: publisher,
		Consumer:  &consumer,
		sessions:  map[string]*melody.Session{},
	}, nil
}

func (c *Client) Close() {
	if err := c.Publisher.Close(); err != nil {
		log.Errorf("error closing Publisher: %s", err)
	}

	if err := c.Consumer.Close(); err != nil {
		log.Errorf("error closing Consumer: %s", err)
	}
}

const routingKey = "routing_key"

func (c *Client) Publish(body []byte, correlationID, exchange string) error {
	err := c.Publisher.Publish(
		body,
		[]string{routingKey},
		rabbitmq.WithPublishOptionsMandatory,
		rabbitmq.WithPublishOptionsPersistentDelivery,
		rabbitmq.WithPublishOptionsExchange(exchange),
		rabbitmq.WithPublishOptionsCorrelationID(correlationID),
	)
	if err != nil {
		log.Errorf("error publishing: %s", err)

		return err
	}

	return nil
}

func (c *Client) Consume(
	consumerName,
	queueName,
	exchange string,
	action func(delivery rabbitmq.Delivery) rabbitmq.Action) error {
	err := c.Consumer.StartConsuming(
		action,
		queueName,
		[]string{routingKey},
		rabbitmq.WithConsumeOptionsConcurrency(1),
		rabbitmq.WithConsumeOptionsQueueDurable,
		rabbitmq.WithConsumeOptionsQuorum,
		rabbitmq.WithConsumeOptionsBindingExchangeName(exchange),
		rabbitmq.WithConsumeOptionsBindingExchangeDurable,
		rabbitmq.WithConsumeOptionsConsumerName(consumerName),
	)
	if err != nil {
		log.Errorf("error starting consuming: %s", err)

		return err
	}

	return nil
}

func (c *Client) StoreSession(key string, session *melody.Session) {
	c.sessions[key] = session
}

func (c *Client) GetSession(key string) *melody.Session {
	return c.sessions[key]
}

func (c *Client) RemoveSession(key string) {
	delete(c.sessions, key)
}

func (c *Client) BroadcastAction(delivery rabbitmq.Delivery) rabbitmq.Action {
	currentSession := c.GetSession(delivery.Delivery.CorrelationId)

	err := c.Websocket.BroadcastFilter([]byte("<system> "+cast.ToString(delivery.Body)), func(s *melody.Session) bool {
		return s.Request.URL.Path == currentSession.Request.URL.Path
	})
	if err != nil {
		log.Errorf("error broadcasting message: %s", err)
	}

	c.RemoveSession(delivery.Delivery.CorrelationId)

	return rabbitmq.Ack
}
