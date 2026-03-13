package tools

import (
	"fmt"

	"github.com/rabbitmq/amqp091-go"
)

var (
	// RabbitMQ 连接字符串
	Rabbitmqch   *amqp091.Channel    //通道
	Rabbitmqconn *amqp091.Connection //连接
	QueueName    = "gochat-messages"
)

func InitRabbitMQ(url string) error {
	var err error

	Rabbitmqconn, err = amqp091.Dial(url)

	if err != nil {
		return fmt.Errorf("连接 RabbitMQ 失败: %v", err)
	}
	Rabbitmqch, err = Rabbitmqconn.Channel()
	if err != nil {
		return fmt.Errorf("创建通道失败: %v", err)
	}
	_, err = Rabbitmqch.QueueDeclare(
		QueueName, // 队列名称
		true,      // 持久化
		false,     // 自动删除
		false,     // 排他性
		false,     // 无等待
		nil,       // 其他属性
	)
	if err != nil {
		return fmt.Errorf("声明队列失败: %v", err)
	}
	return nil
}

func SendMessage(msg []byte) error {
	if Rabbitmqch == nil {
		return fmt.Errorf("RabbitMQ 通道未初始化")
	}
	return Rabbitmqch.Publish(
		"",
		QueueName,
		false,
		false,
		amqp091.Publishing{
			ContentType:  "application/json",
			Body:         msg,
			DeliveryMode: amqp091.Persistent,
		},
	)

}
func ConsumeMessages(callback func([]byte) error) error {
	if Rabbitmqch == nil {
		return fmt.Errorf("rabbitmq channel not init")
	}
	msgs, err := Rabbitmqch.Consume(
		QueueName, // 队列名称
		"",        // 消费者标签
		false,     // 自动确认
		false,     // 排他性
		false,     // 无本地
		false,     // 无等待
		nil,       // 额外参数

	)
	if err != nil {
		return fmt.Errorf("消费者注册失败:%v", err)
	}

	go func() {
		for msg := range msgs {
			if err := callback(msg.Body); err != nil {
				fmt.Printf("处理消息失败:%v\n", err)
			}
			msg.Ack(false)
		}
	}()
	return nil
}

func CloseRabbitMQ() error {
	var err error
	if Rabbitmqch != nil {
		if closeErr := Rabbitmqch.Close(); closeErr != nil {
			err = closeErr
		}
	}
	if Rabbitmqconn != nil {
		if closeErr := Rabbitmqconn.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}
	return err
}
