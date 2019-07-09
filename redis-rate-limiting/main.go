package main

import (
	"log"
	"net"
	"strconv"
	"time"

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

const LUA_SCRIPT_LEAKY_BUCKET = `
		local ret = redis.call("INCRBY", KEYS[1], "1")
		if ret == 1 then
			redis.call("PEXPIRE", KEYS[1], KEYS[2])
		end
		return ret
`

type RateLimiter struct {
	//Token string
	//Take
	l        redis.Conn
	Capacity int
	Rate     int
}

func (rl *RateLimiter) Take(token string) bool {
	redisConn, err := mredigo.Connect(MREDIGO_KEY)
	if err != nil {
		log.Fatal(err)
	}
	defer redisConn.Close()

	luaScript := redis.NewScript(2, LUA_SCRIPT_LEAKY_BUCKET)
	currency, err := redis.Int(luaScript.Do(redisConn, token, strconv.Itoa(rl.Rate)))
	if err != nil {
		log.Fatal(err)
	}

	if currency > rl.Capacity {
		return false
	}

	return true
}

const ()

func main() {
	cfg := mredigo.NewConfig()
	cfg.Host = net.JoinHostPort("172.77.0.99", "6379")
	cfg.Database = "7"

	if err := mredigo.CreatePool(MREDIGO_KEY, true, cfg); err != nil {
		log.Fatal(err)
	}

	rl := &RateLimiter{
		Capacity: 1,
		Rate:     10000,
	}

	start := time.Now()
	for i := 0; i < 1000; i++ {
		for !rl.Take("user/ming") {
		}
		log.Println(time.Now(), i)
	}

	log.Println(time.Now().Sub(start))

	//luaScript := redis.NewScript(1, LUA_SCRIPT_ECHO)
	//s, err := redis.String(luaScript.Do(redisConn, "Hello, world!"))
	//if err != nil {
	//	log.Fatal(err)
	//}

	//fmt.Println(s)
}
