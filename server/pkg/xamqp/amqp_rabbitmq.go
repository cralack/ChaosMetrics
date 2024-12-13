package xamqp

import (
	"context"
	"fmt"
	"time"

	"github.com/cralack/ChaosMetrics/server/internal/global"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type RabbitMQ struct {
	ConnNotify    chan *amqp.Error
	channelNotify chan *amqp.Error
	conn          *amqp.Connection
	Channel       *amqp.Channel
	logger        *zap.Logger
	conf          *Config
	ctx           context.Context
	handler       func([]byte) error
	Delivery      <-chan amqp.Delivery
	role          ROLE
	roleTag       string
}

var _ MessageQueue = &RabbitMQ{}

func NewRabbitMQ(role ROLE, handler func([]byte) error, sets ...Setup) (*RabbitMQ, error) {
	logger := global.ChaLogger
	if global.ChaConf.AmqpConf == nil {
		return nil, fmt.Errorf("AMQP configuration is missing")
	}
	address := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		global.ChaConf.AmqpConf.User,
		global.ChaConf.AmqpConf.Password,
		global.ChaConf.AmqpConf.Host,
		global.ChaConf.AmqpConf.Port,
	)
	conf := &defaultConfig
	conf.Addr = address
	conf.AutoDelete = global.ChaConf.AmqpConf.AutoDelete
	for _, set := range sets {
		set(conf)
	}

	if role == Consumer && handler == nil {
		return nil, fmt.Errorf("no handler for consumer or consumer handler is nil")
	}
	return &RabbitMQ{
		logger:  logger,
		role:    role,
		roleTag: RoleTag[role],
		conf:    conf,
		handler: handler,
		ctx:     conf.Context,
	}, nil
}

func (m *RabbitMQ) Stop() {
	select {
	case <-m.ctx.Done():
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
}

func (m *RabbitMQ) Start() error {
	if err := m.Run(); err != nil {
		return err
	}
	go m.Stop()
	go m.ReConnect()

	return nil
}

func (m *RabbitMQ) Run() (err error) {
	if m.conn, err = amqp.Dial(m.conf.Addr); err != nil {
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
		m.Delivery, err = m.Channel.Consume(
			q.Name,
			m.roleTag,
			false,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			return err
		}
		go m.Handle()
	}

	m.ConnNotify = m.conn.NotifyClose(make(chan *amqp.Error))
	m.channelNotify = m.Channel.NotifyClose(make(chan *amqp.Error))

	m.logger.Info("rabbitmq starting...")
	return
}

func (m *RabbitMQ) ReConnect() {
	for {
		select {
		case err := <-m.ConnNotify:
			if err != nil {
				m.logger.Error("rabbitmq connection NotifyClose", zap.Error(err))
			}
		case err := <-m.channelNotify:
			if err != nil {
				m.logger.Error("rabbitmq Channel NotifyClose", zap.Error(err))
			}
		case <-m.ctx.Done():
			m.logger.Info("rabbitmq quitting...")
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
		for err := range m.ConnNotify {
			m.logger.Error("", zap.Error(err))
		}

	quit: // Reconnect loop
		for {
			select {
			case <-m.ctx.Done():
				m.logger.Debug("quitting reconnect loop")
				return
			default:
				m.logger.Info("rabbitmq reconnecting...")
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

func (m *RabbitMQ) Handle() {
	if m.role == Producer {
		m.logger.Info("producer cannot consume message")
		return
	}
	for d := range m.Delivery {
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

func (m *RabbitMQ) Publish(body []byte, exchange, routingkey string, delay int64) error {
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
		exchange,
		routingkey,
		false,
		false,
		publishing,
	)
}
