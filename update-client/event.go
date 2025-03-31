package main

import (
	"fmt"
	"strings"
	"encoding/json"
	"log"
	"github.com/streadway/amqp"
)


type Consumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	tag     string
	done    chan error
}

type UpdateToken struct {
	Type 				string 	`json:"type"`
	TimeStamp 			int 	`json:"timestamp"`
	OriginService 		string 	`json:"originService"`
	DestinationService 	string 	`json:"destinationService"`
	Id 					string 	`json:"id"`
}

func StartListener(appName string, amqpServer string, exchangeName string) {
	// 监听配置中心的更新事件
	// 1. 监听配置中心的更新事件
	// 2. 接收配置中心的更新事件
	// 3. 更新本地缓存
	// 4. 通知服务端更新配置

	err := NewConsumer(amqpServer, exchangeName, "topic", "springCloudBus", exchangeName, appName)
	if err != nil {
		fmt.Println("Failed to start consumer: ", err)
		return
	}
	select {} // stop a Goroutine from finishing
}

func NewConsumer(amqpURI, exchange, exchangeType, queue, key, ctag string) error {
	c := &Consumer{
		conn: 		nil,
		channel: 	nil,
		tag: 		ctag,
		done: 		make(chan error),
	}
	
	var err error
	// 1. 连接到 RabbitMQ
	c.conn, err = amqp.Dial(amqpURI)
	
	c.channel, err = c.conn.Channel()

	if err = c.channel.ExchangeDeclare(
		exchange,
		exchangeType,
		true,  // durable
		false, // auto-deleted
		false, // internal
		false, // no-wait
		nil,   // arguments
	); err != nil {
		return fmt.Errorf("Failed to declare exchange: %s", err)
	}

	// 声明 channel 中的 Queue
	_, err = c.channel.QueueDeclare(
		queue,
		true,  // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)


	// 绑定消费的队列
	if err = c.channel.QueueBind(
		queue, 
		key,
		exchange,
		false,
		nil,
	); err != nil {
		return fmt.Errorf("Failed to bind queue: %s", err)
	}

	// 消费队列中的消息
	deliveries, err := c.channel.Consume(
		queue,
		c.tag,
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // arguments
	)
	
	go handle(deliveries, c.done)

	return nil
}

func handle(deliveries <-chan amqp.Delivery, done chan error) {
	for d := range deliveries {
		fmt.Printf("Received a message: %s\n", d.Body)
		// 处理消息
		handleRefreshEvent(d.Body, d.ConsumerTag)
		// 处理完消息后，发送 ack
		d.Ack(false)
	}
	log.Printf("handle: deliveries channel closed")
	done <- nil
}

func handleRefreshEvent(body []byte, consumerTag string) {
	// 处理配置中心的更新事件
	// 1. 接收配置中心的更新事件
	// 2. 更新本地缓存
	// 3. 通知服务端更新配置
	// 4. 通知服务端更新配置
	fmt.Println("Received a refresh event")

	updateToken := &UpdateToken{}
	err := json.Unmarshal(body, updateToken)
	if err != nil {
		log.Printf("Failed to unmarshal update token: %v", err.Error())
	} else {
		if strings.Contains(updateToken.DestinationService, "consumerTag") {
			loadRemoteConfig()
		}
	}
}