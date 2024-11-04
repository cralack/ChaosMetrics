package xamqp

import (
	"fmt"
	"time"

	"github.com/cralack/ChaosMetrics/server/internal/config"
	"github.com/cralack/ChaosMetrics/server/internal/global"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type RabbitMQ struct {
	connNotify    chan *amqp.Error
	channelNotify chan *amqp.Error
	conn          *amqp.Connection
	Channel       *amqp.Channel
	logger        *zap.Logger
	conf          *config.AmqpConfig
	done          chan struct{}
	handler       func([]byte) error
	delivery      <-chan amqp.Delivery
	role          ROLE
	roleTag       string
}

var _ global.MessageQueue = &RabbitMQ{}

func NewRabbitMQ(role ROLE, handler func([]byte) error) (*RabbitMQ, error) {
	logger := global.ChaLogger
	conf := global.ChaConf.AmqpConf
	address := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		conf.User, conf.Password, conf.Host, conf.Port)
	if role == Consumer && handler == nil {
		return nil, fmt.Errorf("no handler for consumer or consumer handler is nil")
	}
	return &RabbitMQ{
		logger:  logger,
		role:    role,
		roleTag: RoleTag[role],
		conf: &config.AmqpConfig{
			URL:        address,
			Exchange:   Exchange,
			Queue:      Queue,
			RoutingKey: RoutingKey,
			AutoDelete: conf.AutoDelete,
		},
		handler: handler,
		done:    make(chan struct{}),
	}, nil
}

func (m *RabbitMQ) Stop() {
	close(m.done)

	if !m.conn.IsClosed() {
		if err := m.Channel.Cancel(m.roleTag, true); err != nil {
			m.logger.Error("rabbitmq Channel cancel failed", zap.Error(err))
		}
		if err := m.Channel.Close(); err != nil {
			m.logger.Error("rabbitmq Channel closed failed", zap.Error(err))
		}
	}

	if err := m.conn.Close(); err != nil {
		m.logger.Error("rabbitmq connection close failed", zap.Error(err))
	}
}

func (m *RabbitMQ) Start() error {
	if err := m.Run(); err != nil {
		return err
	}

	go m.ReConnect()

	return nil
}

func (m *RabbitMQ) Run() (err error) {
	if m.conn, err = amqp.Dial(m.conf.URL); err != nil {
		return err
	}
	if m.Channel, err = m.conn.Channel(); err != nil {
		_ = m.conn.Close()
		return err
	}

	// 声明一个主要使用的 exchange
	err = m.Channel.ExchangeDeclare(
		m.conf.Exchange,
		DelayType,
		true,
		m.conf.AutoDelete,
		false,
		false,
		amqp.Table{
			"x-delayed-type": "direct",
		},
	)
	if err != nil {
		return err
	}

	// 声明一个延时队列, 延时消息就是要发送到这里
	q, err := m.Channel.QueueDeclare(
		m.conf.Queue,
		true,
		m.conf.AutoDelete,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	err = m.Channel.QueueBind(
		q.Name,
		m.conf.RoutingKey,
		m.conf.Exchange,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	if m.role == Consumer {
		m.delivery, err = m.Channel.Consume(
			q.Name,
			m.roleTag,
			false,
			false,
			false,
			false,
			nil,
		)
	}

	if err != nil {
		return err
	}

	m.connNotify = m.conn.NotifyClose(make(chan *amqp.Error))
	m.channelNotify = m.Channel.NotifyClose(make(chan *amqp.Error))

	return
}

func (m *RabbitMQ) ReConnect() {
	for {
		select {
		case err := <-m.connNotify:
			if err != nil {
				m.logger.Error("rabbitmq connection NotifyClose", zap.Error(err))
			}
		case err := <-m.channelNotify:
			if err != nil {
				m.logger.Error("rabbitmq Channel NotifyClose", zap.Error(err))
			}
		case <-m.done:
			m.logger.Debug("job done")
			return
		}

		// backstop
		if !m.conn.IsClosed() {
			if err := m.Channel.Cancel(m.roleTag, true); err != nil {
				m.logger.Error("rabbitmq Channel cancel failed", zap.Error(err))
			}
			if err := m.conn.Close(); err != nil {
				m.logger.Error("rabbitmq connection close failed", zap.Error(err))
			}
		}

		// IMPORTANT: HAVE TO Clear notify
		for err := range m.channelNotify {
			m.logger.Error("", zap.Error(err))
		}
		for err := range m.connNotify {
			m.logger.Error("", zap.Error(err))
		}

	quit: // Reconnect loop
		for {
			select {
			case <-m.done:
				return
			default:
				m.logger.Error("rabbitmq reconnecting...")
				if err := m.Run(); err != nil {
					m.logger.Error("rabbitmq connection reconnect failed", zap.Error(err))
					time.Sleep(5 * time.Second)
					continue
				}
				break quit
			}
		}
	}
}

func (m *RabbitMQ) Consume() {
	if m.role == Producer {
		m.logger.Info("producer cannot consume message")
		return
	}
	for d := range m.delivery {
		// m.logger.Debug(string(d.Body))
		go func(delivery amqp.Delivery) {
			if err := m.handler(delivery.Body); err != nil {
				m.logger.Error("rabbitmq handler failed", zap.Error(err))
				_ = delivery.Reject(true)
			} else {
				_ = delivery.Ack(false)
			}
		}(d)
	}
}

func (m *RabbitMQ) Publish(body []byte, delay int64) error {
	publishing := amqp.Publishing{
		ContentType: "text/plain",
		Body:        body,
	}
	if delay >= 0 {
		publishing.Headers = amqp.Table{
			DelayHeader: delay,
		}
	}
	return m.Channel.Publish(
		m.conf.Exchange,
		m.conf.RoutingKey,
		false,
		false,
		publishing,
	)
}
