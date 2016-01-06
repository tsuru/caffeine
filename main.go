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
		app, err := getApp(r.Host)
		if err != nil {
			log.Println(err)
			return
		}

		conn := redisPool.Get()
		defer conn.Close()

		restoreRoute(app.Ip, conn)
		startApp(*app)

		waitBeforeProxy()

		proxy := createProxy(r)
		proxy.ServeHTTP(w, r)
	})

	http.ListenAndServe("0.0.0.0:8888", nil)
}

func restoreRoute(appAddress string, conn redis.Conn) {
	name := HIPACHE_PREFIX + appAddress
	log.Printf("Deleting %s\n", name)
	_, err := conn.Do("LTRIM", name, 0, 0)
	if err != nil {
		log.Printf("Error deleting route: %s\n", err)
	}
}

func waitBeforeProxy() {
	sleepTime, err := strconv.Atoi(getConfig("WAIT_BEFORE_PROXY"))
	if err == nil && sleepTime > 0 {
		log.Printf("Waiting %d seconds before proxying...\n", sleepTime)
		time.Sleep(time.Duration(sleepTime) * time.Second)
	}
}

func hipacheRedisAddr() string {
	host := getConfig("HIPACHE_REDIS_HOST")
	port := getConfig("HIPACHE_REDIS_PORT")

	return fmt.Sprintf("%s:%s", host, port)
}

func hipacheRedisMaxConn() int {
	maxConn, _ := strconv.Atoi(getConfig("HIPACHE_REDIS_MAX_CONN"))
	return maxConn
}

func getConfig(key string) string {
	defaultValues := map[string]string{
		"HIPACHE_REDIS_HOST":     "localhost",
		"HIPACHE_REDIS_PORT":     "6379",
		"HIPACHE_REDIS_MAX_CONN": "10",
		"TSURU_HOST":             "http://localhost",
		"TSURU_APP_PROXY":        "",
		"TSURU_TOKEN":            "",
	}

	value := os.Getenv(key)
	if value == "" {
		return defaultValues[key]
	}

	return value
}
