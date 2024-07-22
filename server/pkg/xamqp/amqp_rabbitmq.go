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
	channel       *amqp.Channel
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

	roleConfig := &config.AmqpConfig{
		Address:    address,
		Exchange:   Constant[role][Exchange],
		Queue:      Constant[role][Queue],
		RoutingKey: Constant[role][RoutingKey],
		AutoDelete: conf.AutoDelete,
	}
	return &RabbitMQ{
		logger:  logger,
		role:    role,
		roleTag: Constant[role][Tag],
		conf:    roleConfig,
		handler: handler,
		done:    make(chan struct{}),
	}, nil
}

func (m *RabbitMQ) Stop() {
	close(m.done)

	if !m.conn.IsClosed() {
		if err := m.channel.Cancel(m.roleTag, true); err != nil {
			m.logger.Error("rabbitmq channel cancel failed", zap.Error(err))
		}
		if err := m.channel.Close(); err != nil {
			m.logger.Error("rabbitmq channel closed failed", zap.Error(err))
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
	if m.conn, err = amqp.Dial(m.conf.Address); err != nil {
		return err
	}
	if m.channel, err = m.conn.Channel(); err != nil {
		_ = m.conn.Close()
		return err
	}

	// 声明一个主要使用的 exchange
	err = m.channel.ExchangeDeclare(m.conf.Exchange, DelayType, true,
		m.conf.AutoDelete, false, false, amqp.Table{"x-delayed-type": "direct"},
	)
	if err != nil {
		return err
	}

	// 声明一个延时队列, 延时消息就是要发送到这里
	q, err := m.channel.QueueDeclare(m.conf.Queue, false, m.conf.AutoDelete, false, false, nil)
	if err != nil {
		return err
	}

	err = m.channel.QueueBind(q.Name, "", m.conf.Exchange, false, nil)
	if err != nil {
		return err
	}

	m.delivery, err = m.channel.Consume(q.Name, m.roleTag, false, false, false, false, nil)
	if err != nil {
		return err
	}

	m.connNotify = m.conn.NotifyClose(make(chan *amqp.Error))
	m.channelNotify = m.conn.NotifyClose(make(chan *amqp.Error))

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
				m.logger.Error("rabbitmq channel NotifyClose", zap.Error(err))
			}
		case <-m.done:
			return
		}

		// backstop
		if !m.conn.IsClosed() {
			if err := m.channel.Cancel(m.roleTag, true); err != nil {
				m.logger.Error("rabbitmq channel cancel failed", zap.Error(err))
			}
			if err := m.conn.Close(); err != nil {
				m.logger.Error("rabbitmq connection close failed", zap.Error(err))
			}
		}

		// Clear notify channels
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
	for d := range m.delivery {
		go func(delivery amqp.Delivery) {
			if err := m.handler(delivery.Body); err != nil {
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
	if delay > 0 {
		publishing.Headers = amqp.Table{
			DelayHeader: delay,
		}
	}
	return m.channel.Publish(m.conf.Exchange, m.conf.RoutingKey, true, false, publishing)
}
