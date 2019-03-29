package util

import (
	"encoding/json"
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
	conn, err := amqp.Dial(GetConfig("rabbit.url"))
	onError(err, "failed to connect tp queue")
	channel, err := conn.Channel()
	onError(err, "failed to open a channel")
	// util
	util = &Util{
		conn:    conn,
		channel: channel,
	}
	return util
}

//push
func HandlePush(info map[string]string) {
	if util == nil {
		util = HandleConn()
	}
	// info
	action := info["action"]
	msg, _ := json.Marshal(info)
	// 定义队列
	queue, err := util.channel.QueueDeclare(
		GetConfig("rabbit.prefix")+action,
		false,
		false,
		false,
		false,
		nil,
	)
	onError(err, "定义队列失败")
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
	onError(err, "发送消息失败")
}

// error
func onError(err error, msg string) {
	if err == nil {
		return
	}
	//if util != nil {
	//	if !util.conn.IsClosed() {
	//		connErr := util.conn.Close()
	//		FailOnErrorNoExit(connErr, "error")
	//	}
	//	chErr := util.channel.Close()
	//	FailOnErrorNoExit(chErr, "error")
	//}
	FailOnErrorNoExit(err, msg)
}
