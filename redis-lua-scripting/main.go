package main

import (
	"fmt"
	"log"
	"net"

	"github.com/gomodule/redigo/redis"
	"github.com/ming-go/pkg/mredigo"
)

const (
	MREDIGO_KEY = "lua-test"
)

const (
	LUA_SCRIPT_ECHO = `
		return redis.call("ECHO", KEYS[1])
	`
)

func main() {
	cfg := mredigo.NewConfig()
	cfg.Host = net.JoinHostPort("172.77.0.99", "6379")
	cfg.Database = "7"

	if err := mredigo.CreatePool(MREDIGO_KEY, true, cfg); err != nil {
		log.Fatal(err)
	}

	redisConn, err := mredigo.Connect(MREDIGO_KEY)
	if err != nil {
		log.Fatal(err)
	}

	luaScript := redis.NewScript(1, LUA_SCRIPT_ECHO)
	s, err := redis.String(luaScript.Do(redisConn, "Hello, world!"))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(s)
}
