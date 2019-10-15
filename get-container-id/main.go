package main

import (
	"encoding/json"
	"errors"
	"flag"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
)

func getConatinerID() (string, error) {
	b, err := ioutil.ReadFile("/proc/1/cpuset")
	if err != nil {
		return "", err
	}

	cpuset := string(b)
	if !strings.Contains(cpuset, "docker") {
		return "", errors.New("Not in docker")
	}

	cpusetSplit := strings.Split(cpuset, "/")
	return cpusetSplit[len(cpusetSplit)-1], nil
}

type stringCache struct {
	s    string
	flag bool
	sync.RWMutex
}

func (sc *stringCache) Get() (string, bool) {
	sc.RLock()
	defer sc.RUnlock()

	return sc.s, sc.flag
}

func (sc *stringCache) Set(s string) {
	sc.Lock()
	sc.s = s
	sc.flag = true
	sc.Unlock()
}

type responseSuccess struct {
	Data string `json:"data"`
}

type errs struct {
	Message string `json:"message"`
}

type responseError struct {
	Errors errs `json:"errors"`
}

var httpPort = "8787"

func main() {
	flag.StringVar(&httpPort, "httpPort", "8787", "-httpPort 8787")
	flag.Parse()

	sc := stringCache{}
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		b, err := json.Marshal(responseSuccess{Data: "Hello, world!"})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	})

	mux.HandleFunc("/containerID", func(w http.ResponseWriter, r *http.Request) {
		containerID, exists := sc.Get()
		if !exists {
			var err error
			containerID, err = getConatinerID()
			if err != nil {
				b, err := json.Marshal(responseError{Errors: errs{Message: err.Error()}})
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				w.Header().Set("Content-Type", "application/json")
				w.Write(b)
				return
			}

			sc.Set(containerID)
		}

		b, err := json.Marshal(responseSuccess{Data: containerID})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	})

	listener, err := net.Listen("tcp", ":"+httpPort)
	if err != nil {
		log.Fatal(err)
	}

	httpServer := &http.Server{
		Handler: mux,
	}

	log.Println("http server started on :" + httpPort)
	httpServer.Serve(listener)
}