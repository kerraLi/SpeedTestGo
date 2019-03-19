package util

import (
	"github.com/streadway/amqp"
)

// 指针
var util *Util

type Util struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func MqUtil() *Util {
	if util == nil {
		util = HandleConn()
	}
	return util
}

// onConn
func HandleConn() *Util {
	var err error
	conn, err := amqp.Dial(GetConfig("queue.url"))
	OnError(err, "failed to connect tp queue")
	channel, err := conn.Channel()
	OnError(err, "failed to open a channel")
	// util
	util = &Util{
		conn:    conn,
		channel: channel,
	}
	return util
}

//push
func HandlePush(action string, msg []byte) {
	if util == nil {
		util = HandleConn()
	}
	// 定义队列
	queue, err := util.channel.QueueDeclare(
		action,
		false,
		false,
		false,
		false,
		nil,
	)
	OnError(err, "定义队列失败")
	// 发送消息
	if len(msg) == 0 {
		msg = []byte("Hello，World！")
	}
	err = util.channel.Publish(
		"",
		queue.Name,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Transient,
			ContentType:  "text/plain",
			Body:         msg,
		})
	OnError(err, "发送消息失败")
}

// error
func OnError(err error, msg string) {
	if err == nil {
		return
	}
	if util != nil {
		if !util.conn.IsClosed() {
			connErr := util.conn.Close()
			FailOnError(connErr, "error")
		}
		chErr := util.channel.Close()
		FailOnError(chErr, "error")
	}
	FailOnError(err, msg)
}
