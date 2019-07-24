package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"
)

const hostPort = "127.0.0.1:8787"

func server() {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-r.Context().Done():
			log.Println("Context Done")
		case <-time.After(7 * time.Minute):
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			io.WriteString(w, "Hello, world!")
		}
	})

	listener, err := net.Listen("tcp", hostPort)
	if err != nil {
		log.Fatal(err)
	}

	httpServer := &http.Server{
		Handler: mux,
	}

	httpServer.Serve(listener)
}

func client() {
	req, err := http.NewRequest("GET", fmt.Sprintf("http://%s%s", hostPort, "/hello"), nil)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(req.Context(), 5*time.Second)
	defer cancel()

	req = req.WithContext(ctx)
	httpClient := http.DefaultClient
	response, err := httpClient.Do(req)
	if err != nil {
		log.Println(err)
		return
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(string(body))
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	go server()
	<-time.After(1 * time.Second)
	start := time.Now()
	client()
	log.Println(time.Now().Sub(start))
}
