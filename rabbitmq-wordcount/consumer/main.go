package main

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	pbWC "github.com/iwdmb/wordcount/producer/pkg/protobuf/wordcount"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

func wordCount(wc *pbWC.WordCount) {
	if wc.Count == nil {
		wc.Count = make(map[string]int64)
	}

	for _, v := range wc.Words {
		if _, ok := wc.Count[v]; ok {
			wc.Count[v] += 1
		} else {
			wc.Count[v] = 1
		}
	}
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}

	ch.Qos(99999, 0, false)

	defer ch.Close()

	msgs, err := ch.Consume(
		"wordcount_queue",
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		zap.L().Info(
			"ch.Consume error",
			zap.Error(err),
		)
	}

	forever := make(chan struct{})

	for d := range msgs {
		wc := pbWC.WordCount{}

		if err := proto.Unmarshal(d.Body, &wc); err != nil {
			zap.L().Info(
				"proto.Unmarshal error",
				zap.Error(err),
			)
		}

		wordCount(&wc)

		fmt.Println(wc.Count)

		//<-time.After(512 * time.Millisecond)
	}

	<-forever
}
