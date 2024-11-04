package test

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/cralack/ChaosMetrics/server/pkg/xamqp"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func Test_messageQue(t *testing.T) {
	// 连接到RabbitMQ服务器
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		logger.Error("Failed to connect to RabbitMQ", zap.Error(err))
	}
	defer func(conn *amqp.Connection) {
		_ = conn.Close()
	}(conn)

	// 创建一个通道
	ch, err := conn.Channel()
	if err != nil {
		logger.Error("Failed to open a channel")
	}
	defer func(ch *amqp.Channel) {
		_ = ch.Close()
	}(ch)

	err = ch.ExchangeDeclare(
		"test_exchange", // 交换机名称
		xamqp.DelayType, // 交换机类型
		true,            // 是否持久化
		true,            // 是否自动删除
		false,           // 是否内部使用
		false,           // 是否等待服务器响应
		amqp.Table{ // 其他属性
			"x-delayed-type": "direct",
		},
	)
	if err != nil {
		logger.Error("Failed to declare an exchange")
	}

	q, err := ch.QueueDeclare(
		"test_queue", // 队列名称，留空表示由RabbitMQ自动生成
		true,         // 是否持久化
		true,         // 是否自动删除（当没有任何消费者连接时）
		false,        // 是否排他队列（仅限于当前连接）
		false,        // 是否等待服务器响应
		nil,
	)
	if err != nil {
		logger.Fatal("Failed to declare a queue")
	}

	err = ch.QueueBind(
		q.Name,            // 队列名称
		"test_routingKey", // 路由键，留空表示接收交换机的所有消息
		"test_exchange",   // 交换机名称
		false,             // 是否等待服务器响应
		nil,               // 其他属性
	)
	if err != nil {
		logger.Error("Failed to bind a queue")
	}

	msgs, err := ch.Consume(
		q.Name, // 队列名称
		"",     // 消费者标识符，留空表示由RabbitMQ自动生成
		false,  // 是否自动应答
		false,  // 是否独占模式（仅限于当前连接）
		false,  // 是否禁止本地消息(RabbitMQ本身不支持)
		false,  // 是否等待服务器响应
		nil,    // 其他属性
	)
	if err != nil {
		logger.Error("Failed to register a consumer")
	}

	timer := time.NewTimer(time.Second * 10)
	go func() {
		for d := range msgs {
			logger.Info("received a message:" + string(d.Body))
			_ = d.Ack(false)
		}
	}()

	time.Sleep(time.Second * 1)

	delayedTime := []int{3000, 5000, 1000}
	for i, t := range delayedTime {
		body := fmt.Sprintf("message%d:%s",
			i, time.Now().Local().Format("2006-01-02 15:04:05"))
		err = ch.Publish(
			"test_exchange",   // 交换机名称
			"test_routingKey", // 路由键，留空表示广播给所有队列
			false,             // 是否等待服务器响应
			false,             // 是否立即将消息写入磁盘
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(body),
				Headers: map[string]interface{}{
					"x-delay": t,
				},
			},
		)
		if err != nil {
			logger.Error("Failed to publish a message")
		}
		logger.Info(fmt.Sprintf("sent message %s", body))
	}

	logger.Info(" [*] Waiting for messages. To exit press CTRL+C")
	<-timer.C
}

func Test_RabbitMQ_Producer(t *testing.T) {

	// 初始化生产者实例
	producer, err := xamqp.NewRabbitMQ(xamqp.Producer, nil) // 生产者不需要处理函数
	assert.NoError(t, err)
	assert.NotNil(t, producer)

	err = producer.Start()
	assert.NoError(t, err)

	// should be 4,5,1,2,3
	delayedTime := []int64{9, 1, 6, 4, 3, 8, 2, 5, 7, 10, 0}
	for _, tim := range delayedTime {
		massage := fmt.Sprintf("%d test message", tim)
		err = producer.Publish([]byte(massage), tim*500)
		assert.NoError(t, err)
		if err == nil {
			logger.Debug("sent message " + massage)
		}
	}

	timer := time.NewTimer(time.Second * 15)
	<-timer.C
	producer.Stop()
}

func Test_RabbitMQ_Consumer(t *testing.T) {
	// 初始化消费者实例
	// consumer, err := xamqp.NewRabbitMQ(xamqp.Consumer, nil)
	consumer, err := xamqp.NewRabbitMQ(xamqp.Consumer, mockHandler())
	assert.NoError(t, err)
	assert.NotNil(t, consumer)

	err = consumer.Start()
	assert.NoError(t, err)
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)

	// 确保消费者已启动
	go consumer.Consume()
	logger.Info(" [*] Waiting for messages. To exit press CTRL+C")
	<-exit
	consumer.Stop()
}

func mockHandler() func(body []byte) error {
	return func(body []byte) error {
		logger.Info("Received message: " + string(body))
		return nil
	}
}
