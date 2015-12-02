package main

import (
	"log"
	"net/http"
	"net/http/httputil"
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
	conn := redisPool.Get()
	defer conn.Close()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		restoreRoute(r, conn)
		startApp(r)
		time.Sleep(10 * time.Second)

		director := func(req *http.Request) {
			req.URL = r.URL
			req.URL.Scheme = "http"
			req.URL.Host = r.Host
			req.Host = r.Host
			log.Printf("Request: %#v", req)
		}

		proxy := &httputil.ReverseProxy{Director: director}
		proxy.ServeHTTP(w, r)
	})

	http.ListenAndServe("0.0.0.0:8888", nil)
}

func restoreRoute(r *http.Request, conn redis.Conn) {

	name := HIPACHE_PREFIX + r.Host
	log.Print("Deleting ", name)
	_, err := conn.Do("LTRIM", name, 0, 0)
	if err != nil {
		log.Printf("Err: %s", err)
	}
}
