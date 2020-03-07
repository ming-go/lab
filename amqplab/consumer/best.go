package main

import (
	"log"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gammazero/deque"
	jsoniter "github.com/json-iterator/go"
	"github.com/streadway/amqp"
)

var configInput amqp.URI
var configOutput amqp.URI

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func init() {
	configInput = amqp.URI{
		Scheme:   "amqp",
		Host:     "192.168.7.102",
		Port:     8787,
		Username: "guest",
		Password: "guest",
		Vhost:    "/",
	}

	configOutput = amqp.URI{
		Scheme:   "amqp",
		Host:     "192.168.7.102",
		Port:     8888,
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

//var json = jsoniter.ConfigCompatibleWithStandardLibrary

func open(config amqp.URI, queueName string) (*amqp.Connection, *amqp.Channel, error) {
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

// LWrapper -
type LWrapper struct {
	//l *list.List
	q deque.Deque
	sync.RWMutex
}

// PushFront -
func (lw *LWrapper) PushFront(v interface{}) {
	lw.Lock()
	//lw.l.PushFront(v)
	lw.q.PushFront(v)
	lw.Unlock()
}

// PopBack -
func (lw *LWrapper) PopBack() interface{} {
	lw.Lock()
	defer lw.Unlock()
	if lw.q.Len() != 0 {
		return lw.q.PopBack()
	}

	return nil
	//item := lw.l.Back()
	//if item != nil {
	//	return lw.l.Remove(item)
	//}

	//return nil
}

// Len -
func (lw *LWrapper) Len() int {
	lw.RLock()
	defer lw.RUnlock()
	return lw.q.Len()
	//return lw.l.Len()
}

// NewLWrapper -
func NewLWrapper() LWrapper {
	return LWrapper{
		//l: list.New(),
		q: deque.Deque{},
	}
}

func main() {
	lw := NewLWrapper()

	// In
	go func() {
		_, ch, err := open(configInput, inName)
		if err != nil {
			log.Fatal(err)
		}

		chDelivery, err := ch.Consume(
			inName,
			"",
			false,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			log.Fatal(err)
		}

		for d := range chDelivery {
			d.Ack(false)
			var data Data
			json.Unmarshal(d.Body, &data)
			data.From = "Consumer"
			lw.PushFront(data)
		}
	}()

	//go func() {
	//	//for {
	//	//	<-time.After(60 * time.Second)
	//	//	runtime.GC()
	//	//}
	//}()

	// Out
	go func() {
		chs := make([]*amqp.Channel, 0, 64)

		for i := 0; i < 64; i++ {
			_, ch, err := open(configOutput, outName[i%4])
			chs = append(chs, ch)
			if err != nil {
				log.Println(err)
			}
		}

		//_, ch, err := open(configOutput, outName[0])
		//if err != nil {
		//	log.Fatal(err)
		//}

		var count uint64

		go func() {
			for {
				log.Println(atomic.LoadUint64(&count) / 4)
				<-time.After(1 * time.Second)
			}
		}()

		for i := 0; i < 6; i++ {
			go func() {
				for {
					if lw.Len() != 0 {
						var data interface{}
						data = lw.PopBack()

						b, _ := json.Marshal(data)
						for i := 0; i < 4; i++ {
							chs[atomic.AddUint64(&count, 1)%64].Publish(

								"",
								outName[i],
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
				}
			}()
		}
	}()

	forever := make(chan struct{})
	<-forever
}
