package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	consulAPI "github.com/hashicorp/consul/api"
	"github.com/ming-go/pkg/mhttp"
)

var consulAddress = ""

const SERVICE_NAME = "service-discovery-lab"

var httpHost = ""
var httpPort = 8787

func main() {
	if len(os.Args) < 6 {
		log.Fatal("you must specify a httpHost & httpPort & consul address")
	}

	flag.StringVar(&httpHost, "host", "127.0.0.1", "http host")
	flag.IntVar(&httpPort, "port", 8787, "http port")
	flag.StringVar(&consulAddress, "consul", "127.0.0.1:8500", "consul addr")
	flag.Parse()

	consulDefaultConfig := consulAPI.DefaultConfig()
	consulDefaultConfig.Address = consulAddress

	client, err := consulAPI.NewClient(consulDefaultConfig)
	if err != nil {
		log.Fatal(err)
	}

	AGENT_SERVICE_ID := fmt.Sprintf("%s-%s:%d", SERVICE_NAME, httpHost, httpPort)

	agent := client.Agent()
	if err := agent.ServiceRegister(&consulAPI.AgentServiceRegistration{
		ID:      AGENT_SERVICE_ID,
		Name:    SERVICE_NAME,
		Address: httpHost,
		Port:    httpPort,
		Check: &consulAPI.AgentServiceCheck{
			CheckID: AGENT_SERVICE_ID,
			TTL:     time.Duration(5 * time.Second).String(),
		},
	}); err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/sd/info", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Server Info\n"))
		w.Write([]byte(fmt.Sprintf("IP: %s\n", httpHost)))
		w.Write([]byte(fmt.Sprintf("Port: %d", httpPort)))
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, world!"))
	})
	mux.HandleFunc("/sd/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	cfg := mhttp.NewConfig().SetMux(mux).SetAddr(httpHost, strconv.Itoa(httpPort))
	srv := mhttp.NewDefaultServer(cfg)

	// Run a HTTP Server
	go func() {
		log.Printf("http server started on %s:%d\n", httpHost, httpPort)
		srv.Run()
	}()

	// Monitor & HealthCheck
	go func() {
		const STATUS_PASS = "pass"
		const STATUS_FAIL = "fail"
		for {
			err := pingServer(2)
			if err != nil {
				agent.UpdateTTL(AGENT_SERVICE_ID, err.Error(), STATUS_FAIL)
			} else {
				agent.UpdateTTL(AGENT_SERVICE_ID, "", STATUS_PASS)
			}
			<-time.After(1 * time.Second)
		}
	}()

	idleConnsClosed := make(chan struct{})

	// Graceful shutdown
	go func() {
		sigint := make(chan os.Signal, 1)

		signal.Notify(sigint, os.Interrupt)
		signal.Notify(sigint, syscall.SIGTERM)

		<-sigint

		srv.Shutdown()
		err := agent.ServiceDeregister(AGENT_SERVICE_ID)
		if err != nil {
			log.Println("agent.ServiceDeregister: %v", err)
		}

		log.Println("Graceful shutdown")
		defer close(idleConnsClosed)
	}()

	<-idleConnsClosed
}

func pingServer(retries int) error {
RETRY:
	resp, err := http.Get(fmt.Sprintf("http://%s/sd/health", net.JoinHostPort(httpHost, strconv.Itoa(httpPort))))
	if err == nil && resp.StatusCode == 200 {
		resp.Body.Close()
		return nil
	}

	if resp != nil {
		resp.Body.Close()
	}

	if retries > 0 {
		<-time.After(1 * time.Second)
		retries--
		goto RETRY
	}

	return errors.New("Cannot connect to the router.")
}
