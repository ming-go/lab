package main

import (
	"log"
	"net"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/streadway/amqp"
)

var config amqp.URI

func init() {
	config = amqp.URI{
		Scheme:   "amqp",
		Host:     "172.77.0.88",
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

func main() {
	// In
	_, ch, err := open(inName)
	if err != nil {
		log.Fatal(err)
	}

	chDelivery, err := ch.Consume(
		outName[3],
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

	for range chDelivery {
	}
}
