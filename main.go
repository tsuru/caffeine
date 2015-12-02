package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/garyburd/redigo/redis"
)

const HIPACHE_PREFIX = "frontend:"

var (
	maxConn, _ = strconv.Atoi(os.Getenv("REDIS_MAX_CONN"))
	redisPool  = redis.NewPool(func() (redis.Conn, error) {
		c, err := redis.Dial("tcp", os.Getenv("REDIS_ADDR"))
		if err != nil {
			return nil, err
		}
		return c, err
	}, maxConn)
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		conn := redisPool.Get()
		defer conn.Close()
		host := r.Host
		restoreRoute(host, conn)
		startApp(host)
		time.Sleep(10 * time.Second)
		proxy := createProxy(r)
		proxy.ServeHTTP(w, r)
	})

	http.ListenAndServe("0.0.0.0:8888", nil)
}

func restoreRoute(host string, conn redis.Conn) {
	name := HIPACHE_PREFIX + host
	log.Print("Deleting ", name)
	_, err := conn.Do("LTRIM", name, 0, 0)
	if err != nil {
		log.Printf("Err: %s", err)
	}
}
