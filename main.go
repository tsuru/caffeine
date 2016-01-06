package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/garyburd/redigo/redis"
)

const HIPACHE_PREFIX = "frontend:"

var (
	redisPool = redis.NewPool(func() (redis.Conn, error) {
		c, err := redis.Dial("tcp", hipacheRedisAddr())
		if err != nil {
			return nil, err
		}
		return c, err
	}, hipacheRedisMaxConn())
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
	log.Printf("Deleting %s\n", name)
	_, err := conn.Do("LTRIM", name, 0, 0)
	if err != nil {
		log.Printf("Err: %s\n", err)
	}
}

func hipacheRedisAddr() string {
	host := os.Getenv("HIPACHE_REDIS_HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("HIPACHE_REDIS_PORT")
	if port == "" {
		port = "6379"
	}

	return fmt.Sprintf("%s:%s", host, port)
}

func hipacheRedisMaxConn() int {
	maxConnValue := os.Getenv("HIPACHE_REDIS_MAX_CONN")
	if maxConnValue == "" {
		maxConnValue = "10"
	}
	maxConn, _ := strconv.Atoi(maxConnValue)
	return maxConn
}
