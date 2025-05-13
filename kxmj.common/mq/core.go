package mq

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"kxmj.common/log"
	"sync"
)

type RabbitmqConfig struct {
	Addr     string `yaml:"addr"`     //地址
	Port     int    `yaml:"port"`     //端口
	User     string `yaml:"user"`     //用户
	Password string `yaml:"password"` //密码
}

type RabbitService struct {
	conn   *amqp.Connection
	close  chan struct{}
	config *RabbitmqConfig
	locker sync.Mutex
}

var rabbitMq *RabbitService

func create(config *RabbitmqConfig) *RabbitService {
	service := RabbitService{
		close:  make(chan struct{}),
		config: config,
		locker: sync.Mutex{},
	}

	conn, err := service.getConnection()
	if err != nil {
		fmt.Println(fmt.Sprintf("can not create rabbitmq connection err:%v", err))
		panic(err)
	}

	service.conn = conn
	return &service
}

func Init(config *RabbitmqConfig) {
	rabbitMq = create(config)
}

func Default() *RabbitService {
	return rabbitMq
}

func (r *RabbitService) getConnection() (conn *amqp.Connection, err error) {
	defer r.locker.Unlock()
	r.locker.Lock()

	url := fmt.Sprintf("amqp://%s:%s@%s:%d", r.config.User, r.config.Password, r.config.Addr, r.config.Port)

	//// RabbitMQ分配的用户名称
	//var user string = "admin"
	//// RabbitMQ用户的密码
	//var pwd string = "bojiu123"
	//// RabbitMQ Broker 的ip地址
	//var host string = "192.168.0.188"
	//// RabbitMQ Broker 监听的端口
	//var port string = "5672"
	//url := "amqp://" + user + ":" + pwd + "@" + host + ":" + port + "/"

	return amqp.Dial(url)
}

func (r *RabbitService) Publish(routerKey string, msg interface{}) error {
	channel, err := r.conn.Channel()
	if err != nil {
		// 断线重连
		conn, err := r.getConnection()
		if err != nil {
			return err
		}

		r.conn = conn
		return r.Publish(routerKey, msg)
	}

	defer channel.Close()

	q, err := channel.QueueDeclare(
		// 队列名称
		routerKey, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)

	if err != nil {
		log.Sugar().Error(fmt.Sprintf("QueueDeclare err:%v", err))
		return err
	}

	body, err := json.Marshal(msg)
	if err != nil {
		log.Sugar().Error(fmt.Sprintf("Marshal err:%v", err))
	}

	err = channel.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			// 将消息标记为持久消息
			DeliveryMode: amqp.Persistent,
			Body:         body,
		})

	return err
}

func (r *RabbitService) Receiver(routerKey string, receive func(msg string) error) {
	channel, err := r.conn.Channel()
	if err != nil {
		log.Sugar().Error(fmt.Sprintf("can not create rabbitmq channel err:%v", err))
		return
	}

	defer channel.Close()

	q, err := channel.QueueDeclare(
		routerKey, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)

	if err != nil {
		log.Sugar().Error(fmt.Sprintf("create rabbitmq queue err:%v", err))
		return
	}

	// 将预取计数器设置为1
	// 在并行处理中将消息分配给不同的工作进程
	err = channel.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)

	if err != nil {
		log.Sugar().Error(fmt.Sprintf("create rabbitmq worker thread err:%v", err))
		return
	}

	msgList, err := channel.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)

	if err != nil {
		log.Sugar().Error(fmt.Sprintf("create rabbitmq consume err:%v", err))
		return
	}

	forever := make(chan bool)
	for d := range msgList {
		err = receive(string(d.Body))
		if err != nil {
			log.Sugar().Error(fmt.Sprintf("process msg err:%v", err))
			err = d.Ack(true)
		} else {
			err = d.Ack(false)
		}

		if err != nil {
			log.Sugar().Error(fmt.Sprintf("ack rabbit queue err:%v", err))
		}
	}
	<-forever
}

func (r *RabbitService) Stop() error {
	return r.conn.Close()
}
