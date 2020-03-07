package main

import (
	"log"
	"net"
	"sync/atomic"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/streadway/amqp"
)

var config amqp.URI

func init() {
	config = amqp.URI{
		Scheme:   "amqp",
		Host:     "172.77.0.87",
		Port:     5672,
		Username: "guest",
		Password: "guest",
		Vhost:    "/",
	}
}

const (
	queueName = "producer"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func open() (*amqp.Connection, *amqp.Channel, error) {
	conn, err := amqp.DialConfig(config.String(), amqp.Config{
		Dial: func(network, address string) (net.Conn, error) {
			return net.DialTimeout(network, address, 1*time.Second)
		},
	})
	if err != nil {
		return nil, nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, nil, err
	}

	if _, err := channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // autoDelete
		false,     // exclusive
		false,     // noWait
		nil,       // args
	); err != nil {
		channel.Close()
		conn.Close()
		return nil, nil, err
	}

	return conn, channel, nil
}

// Data data
type Data struct {
	Username  string
	CreatedAt time.Time
}

const (
	contentTypeApplicationJSON = "application/json"
)

var (
	// N n
	N = 10000000
)

func main() {
	_, ch, err := open()
	if err != nil {
		log.Fatal(err)
	}

	var count uint64

	go func() {
		for {
			<-time.After(1 * time.Second)
			log.Println(atomic.CompareAndSwapUint64(&count, count, 0))
		}
	}()

	for i := 0; i < N; i++ {
		if atomic.LoadUint64(&count) > 10000 {
			i--
			continue
		}

		b, err := json.Marshal(
			Data{
				Username:  "ming",
				CreatedAt: time.Now(),
			},
		)
		if err != nil {
			log.Println(err)
		}

		err = ch.Publish(
			"",
			queueName,
			false,
			false,
			amqp.Publishing{
				DeliveryMode: amqp.Persistent,
				ContentType:  contentTypeApplicationJSON,
				Body:         b,
			},
		)

		atomic.AddUint64(&count, 1)

		if err != nil {
			log.Println(err)
		}
	}

	<-time.After(30 * time.Second)

}
