package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		app, err := getApp(r.Host)
		if err != nil {
			log.Println(err)
			return
		}

		startApp(*app)

		waitBeforeProxy(nil)

		proxy := createProxy(r)
		proxy.ServeHTTP(w, r)
	})
	http.ListenAndServe("0.0.0.0:8888", nil)
}

func getConfig(key string) string {
	defaultValues := map[string]string{
		"TSURU_HOST":        "http://localhost",
		"TSURU_APP_PROXY":   "",
		"TSURU_TOKEN":       "",
		"WAIT_BEFORE_PROXY": "0",
	}

	value := os.Getenv(key)
	if value == "" {
		return defaultValues[key]
	}

	return value
}

func waitBeforeProxy(sleep func(time.Duration)) {
	if sleep == nil {
		sleep = time.Sleep
	}
	sleepTime, err := strconv.Atoi(getConfig("WAIT_BEFORE_PROXY"))
	if err == nil && sleepTime > 0 {
		log.Printf("Waiting %d seconds before proxying...\n", sleepTime)
		sleep(time.Duration(sleepTime) * time.Second)
	}
}
