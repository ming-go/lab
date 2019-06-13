package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	consulAPI "github.com/hashicorp/consul/api"
)

const consulAddress = "172.77.0.22:8500"

var httpHost = ""
var httpPort = 8787

func main() {
	if len(os.Args) < 4 {
		log.Fatal("you must specify a httpHost & httpPort")
	}

	flag.StringVar(&httpHost, "host", "127.0.0.1", "http host")
	flag.IntVar(&httpPort, "port", 8787, "http port")
	flag.Parse()

	consulDefaultConfig := consulAPI.DefaultConfig()
	consulDefaultConfig.Address = consulAddress

	client, err := consulAPI.NewClient(consulDefaultConfig)
	if err != nil {
		log.Fatal(err)
	}

	const SERVICE_NAME = "service-discovery-test"

	agent := client.Agent()
	if err := agent.ServiceRegister(&consulAPI.AgentServiceRegistration{
		ID:      SERVICE_NAME + httpHost + ":" + strconv.Itoa(httpPort),
		Name:    SERVICE_NAME,
		Address: httpHost,
		Port:    httpPort,
		//Check: consulAPI.AgentServiceCheck{
		//	TTL: time.Duration(10 * time.Second).String(),
		//},
	}); err != nil {
		log.Fatal(err)
	}

	//agent.UpdateTTL()
	mux := http.NewServeMux()
	mux.HandleFunc("/sd/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	if err := http.ListenAndServe(net.JoinHostPort(httpHost, strconv.Itoa(httpPort)), nil); err != nil {
		log.Fatal(err)
	}
}

func pingServer(retries int) error {
RETRY:
	resp, err := http.Get(net.JoinHostPort(httpHost, strconv.Itoa(httpPort)))
	if err == nil && resp.StatusCode == 200 {
		return nil
	}

	if retries > 0 {
		<-time.After(1 * time.Second)
		retries--
		goto RETRY
	}

	return nil
}
