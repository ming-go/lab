package main

import (
	"io/ioutil"
	"strings"

	"github.com/golang/protobuf/proto"
	pbWC "github.com/iwdmb/wordcount/producer/pkg/protobuf/wordcount"
	"github.com/ming-go/pkg/util"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

func main() {
	l, err := zap.NewDevelopment()
	zap.ReplaceGlobals(l)

	b, err := ioutil.ReadFile("wordcount.txt")
	if err != nil {
		panic(err)
	}

	w := strings.Replace(string(b), "\n", " ", -1)
	wc := strings.Split(w, " ")

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"wordcount_queue",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}

	splits := util.IntervalSplitBySize(0, int64(len(wc)-1), 2)
	for _, v := range splits {
		b, _ := proto.Marshal(&pbWC.WordCount{
			Words: wc[v[0] : v[1]+1],
		})

		p := amqp.Publishing{
			ContentType: "application/protobuf",
			Body:        b,
		}

		err := ch.Publish("", q.Name, false, false, p)
		if err != nil {
			zap.L().Info(
				"ch.Publish error",
				zap.Error(err),
			)
		}
	}

	//
	/*

		fmt.Println()
	*/

	/*
	 */
}
