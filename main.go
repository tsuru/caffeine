package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	customHeaderValue, _ := getConfig("CUSTOM_HEADER_VALUE")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		app, err := getApp(r.Host)
		if err != nil {
			log.Println(err)
			return
		}

		startApp(*app)

		waitBeforeProxy(nil)

		proxy := createProxy(r, customHeaderValue)
		proxy.ServeHTTP(w, r)
	})
	http.ListenAndServe("0.0.0.0:8888", nil)
}

func getConfig(key string) (string, error) {
	defaultValues := map[string]string{
		"CUSTOM_HEADER_VALUE": "",
		"TSURU_HOST":          "http://localhost",
		"TSURU_APP_PROXY":     "",
		"TSURU_TOKEN":         "",
		"WAIT_BEFORE_PROXY":   "0",
	}

	value := os.Getenv(key)
	if value == "" {
		value = defaultValues[key]
	}

	if value == "" {
		return value, fmt.Errorf("Error, environment variable %s is not defined", key)
	}
	return value, nil
}

func waitBeforeProxy(sleep func(time.Duration)) {
	if sleep == nil {
		sleep = time.Sleep
	}
	waitBeforeProxy, _ := getConfig("WAIT_BEFORE_PROXY")
	sleepTime, err := strconv.Atoi(waitBeforeProxy)
	if err == nil && sleepTime > 0 {
		log.Printf("Waiting %d seconds before proxying...\n", sleepTime)
		sleep(time.Duration(sleepTime) * time.Second)
	}
}
