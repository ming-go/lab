package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"

	consulAPI "github.com/hashicorp/consul/api"
)

const consulAddress = "172.77.0.66:8500"

func main() {
	if len(os.Args) < 4 {
		log.Fatal("you must specify a httpHost & httpPort")
	}

	httpHost := ""
	httpPort := 8787

	flag.StringVar(&httpHost, "host", "127.0.0.1", "http host")
	flag.IntVar(&httpPort, "port", 8787, "http port")
	flag.Parse()

	consulDefaultConfig := consulAPI.DefaultConfig()
	consulDefaultConfig.Address = consulAddress

	client, err := consulAPI.NewClient(consulDefaultConfig)
	if err != nil {
		log.Fatal(err)
	}

	agent := client.Agent()
	if err := agent.ServiceRegister(&consulAPI.AgentServiceRegistration{
		Name:    "service-discovery-test",
		Address: httpHost,
		Port:    httpPort,
		//Check: consulAPI.AgentServiceCheck{
		//	TTL: time.Duration(10 * time.Second).String(),
		//},
	}); err != nil {
		log.Fatal(err)
	}

	//agent.UpdateTTL()

	if err := http.ListenAndServe(net.JoinHostPort(httpHost, strconv.Itoa(httpPort)), nil); err != nil {
		log.Fatal(err)
	}
}
