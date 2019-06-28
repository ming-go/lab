package main

import (
	"encoding/json"
	"flag"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/hashicorp/raft"
	"github.com/ming-go/lab/hashicorp-raft/hraft"
	"github.com/ming-go/pkg/mhttp"
)

var listenHost = ""
var listenPort = 8787
var raftHost = ""
var raftPort = 8787
var nodeID = ""
var master = false

func main() {
	flag.StringVar(&raftHost, "host", "127.0.0.1", "host")
	flag.IntVar(&raftPort, "port", 8787, "port")
	flag.StringVar(&nodeID, "nodeid", "nodeid-01", "nodeid")
	flag.BoolVar(&master, "master", false, "master")

	flag.Parse()

	raftBindAddr := net.JoinHostPort(raftHost, strconv.Itoa(raftPort))
	nodeID = raftBindAddr

	_, err := hraft.New(raft.ServerID(nodeID), raftBindAddr, master)
	if err != nil {
		log.Println(err)
	}

	mux := http.NewServeMux()
	mux.Handle("join", func(w http.ResponseWriter, r *http.Request) {
		type Request struct {
			RaftBindAddr `json:"raftBindAddr"`
		}

		decoder := json.NewDecoder(req.Body)

		req := new(Request)
		err := decoder.Decode(&req)
		if err != nil {
			log.Println(err)
		}

		err := hraft.Join(req.RaftBindAddr, req.RaftBindAddr)
		if err != nil {
			log.Println(err)
		}
	})

	if master() {
		cfg := mhttp.NewConfig().SetMux(mux).SetAddr("127.0.0.1:8788")
		s := mhttp.NewDefaultServer(cfg)
		s.Run()
	}

	<-time.After(10 * time.Minute)
}
