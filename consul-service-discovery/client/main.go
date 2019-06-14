package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	consulAPI "github.com/hashicorp/consul/api"
)

const SERVICE_NAME = "service-discovery-lab"

func roundRobin(lists []string) func() string {
	cRRTarget := -1 // c means closures
	cListsLength := len(lists)
	return func() string {
		cRRTarget = ((cRRTarget + 1) % cListsLength)
		return lists[cRRTarget]
	}
}

type healthwrapper struct {
	health *consulAPI.Health
}

func NewHealthWrapper(h *consulAPI.Health) *healthwrapper {
	return &healthwrapper{
		health: h,
	}
}

func (hw *healthwrapper) GetPassingList() ([]string, error) {
	serviceEntry, _, err := hw.health.Service(SERVICE_NAME, "", true, nil)
	if err != nil {
		return nil, err
	}

	services := make([]string, len(serviceEntry))

	for i, v := range serviceEntry {
		services[i] = fmt.Sprintf("%s:%d", v.Service.Address, v.Service.Port)
	}

	return services, nil
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("you must specify a consul address")
	}

	consulAddress := ""

	flag.StringVar(&consulAddress, "consul", "127.0.0.1:8500", "consul addr")
	flag.Parse()

	consulDefaultConfig := consulAPI.DefaultConfig()
	consulDefaultConfig.Address = consulAddress

	client, err := consulAPI.NewClient(consulDefaultConfig)
	if err != nil {
		log.Fatal(err)
	}

	health := client.Health()

	hw := &healthwrapper{
		health: health,
	}

	services, err := hw.GetPassingList()
	if err != nil {
		log.Fatal(err)
	}

	nextTarget := roundRobin(services)
	for {
		resp, err := http.Get("http://" + nextTarget() + "/sd/info")
		if err != nil {
			log.Println(err)
		} else {
			b, _ := ioutil.ReadAll(resp.Body)
			fmt.Println(string(b))
			fmt.Println()
			resp.Body.Close()
		}

		<-time.After(1 * time.Second)
	}
}
