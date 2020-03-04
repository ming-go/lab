package main

import (
	"log"
	"net"
	"sync"
	"time"

	"container/list"

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

var (
	inName  = "producer"
	outName = []string{
		"out_queue_1",
		"out_queue_2",
		"out_queue_3",
		"out_queue_4",
	}
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func open(queueName string) (*amqp.Connection, *amqp.Channel, error) {
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
	From      string
	CreatedAt time.Time
}

const (
	contentTypeApplicationJSON = "application/json"
)

var (
	// N n
	N = 100000
)

type LWrapper struct {
	l *list.List
	sync.RWMutex
}

func NewLWrapper() LWrapper {
	return LWrapper{
		l: list.New(),
	}
}

func main() {
	lw := NewLWrapper()

	// In
	go func() {
		_, ch, err := open(inName)
		if err != nil {
			log.Fatal(err)
		}

		chDelivery, err := ch.Consume(
			inName,
			"",
			true,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			log.Fatal(err)
		}

		for d := range chDelivery {
			var data Data
			json.Unmarshal(d.Body, &data)
			data.From = "Consumer"
			lw.Lock()
			lw.l.PushFront(data)
			lw.Unlock()
		}
	}()

	// Out
	go func() {
		_, ch, err := open(outName[0])
		if err != nil {
			log.Fatal(err)
		}

		for {
			if lw.l.Len() != 0 {
				lw.Lock()
				item := lw.l.Back()
				var data interface{}
				if item != nil {
					data = lw.l.Remove(item)
				}
				lw.Unlock()

				b, _ := json.Marshal(data)

				ch.Publish(
					"",
					outName[0],
					false,
					false,
					amqp.Publishing{
						DeliveryMode: amqp.Persistent,
						ContentType:  contentTypeApplicationJSON,
						Body:         b,
					},
				)
			}
		}
	}()

	forever := make(chan struct{})
	<-forever
}
