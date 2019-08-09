package main

import (
	"io"
	"net/http"
	"strings"

	"github.com/ming-go/pkg/mhttp"
)

var httpAddr = "127.0.0.1:1024"

func corsWrapper(allowOrigins []string, w http.ResponseWriter) func(origin string) {
	originMap := make(map[string]struct{}, len(allowOrigins))

	for _, origin := range allowOrigins {
		originMap[origin] = struct{}{}
	}

	return func(origin string) {
		_, ok := originMap[origin]
		if ok {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}
	}
}

func server() {
	mux := http.NewServeMux()
	mux.HandleFunc("/cors", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "https://www.test-cors.org")
		w.Header().Set("Access-Control-Allow-Methods", "GET")
		w.Header().Set("Access-Control-Allow-Headers", "test")

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "Hello, world!")
	})

	cfg := mhttp.NewConfig().SetMux(mux).SetAddr(strings.Split(httpAddr, ":")[0], strings.Split(httpAddr, ":")[1])
	s := mhttp.NewDefaultServer(cfg)
	s.Run()
}

func main() {
	go server()
	forever := make(chan struct{})
	<-forever
}
