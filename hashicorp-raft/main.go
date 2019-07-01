package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	hhttp "github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/raft"
	"github.com/ming-go/lab/hashicorp-raft/hraft"
	"github.com/ming-go/pkg/mhttp"
	"github.com/ming-go/pkg/mstring"
)

var joinAddr = ""
var httpAddr = "127.0.0.1:8080"
var raftAddr = "127.0.0.1:8787"

type PostJoinRequest struct {
	RaftAddr string `json:"RaftAddr"`
}

func main() {
	flag.StringVar(&raftAddr, "raftAddr", "127.0.0.1:8787", "-raftAddr 127.0.0.1:8787")
	flag.StringVar(&httpAddr, "httpAddr", "127.0.0.1:8080", "-httpAddr 127.0.0.1:8080")
	flag.StringVar(&joinAddr, "joinAddr", "", "-joinAddr 127.0.0.1:8080")

	flag.Parse()

	nodeID := raftAddr
	nodeIsMaster := mstring.IsBlank(joinAddr)

	hraft, err := hraft.New(raft.ServerID(nodeID), raftAddr, nodeIsMaster)
	if err != nil {
		log.Println(err)
	}

	// master
	if nodeIsMaster {
		mux := http.NewServeMux()
		mux.HandleFunc("/join", func(w http.ResponseWriter, r *http.Request) {
			req := new(PostJoinRequest)
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				log.Println(err)
				return
			}

			err = hraft.Join(req.RaftAddr, req.RaftAddr)
			if err != nil {
				log.Println(err)
				w.Header().Set("Content-Type", "text/plain; charset=utf-8")
				w.WriteHeader(http.StatusBadRequest)
				io.WriteString(w, "Join fail :(")
				return
			}

			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			io.WriteString(w, "Hello, world!")
		})

		cfg := mhttp.NewConfig().SetMux(mux).SetAddr(strings.Split(httpAddr, ":")[0], strings.Split(httpAddr, ":")[1])
		s := mhttp.NewDefaultServer(cfg)
		s.Run()
		return
	}

	// slave
	request := PostJoinRequest{
		RaftAddr: raftAddr,
	}

	log.Println(request.RaftAddr)

	b, err := json.Marshal(request)
	if err != nil {
		log.Println(err)
		return
	}

	req, err := http.NewRequest("POST", joinAddr+"/join", bytes.NewBuffer(b))
	if err != nil {
		log.Println(err)
		return
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")

	resp, err := hhttp.DefaultClient().Do(req)
	if err != nil {
		log.Println(err)
		return
	}

	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	log.Println(string(bodyBytes))

	log.Println(hraft.Raft.Leader())

	<-time.After(10 * time.Minute)
}
