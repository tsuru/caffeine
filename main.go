package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/garyburd/redigo/redis"
)

const HIPACHE_PREFIX = "frontend:"

type App struct {
	Units []Unit
}

type Unit struct {
	ProcessName string
	AppName     string
	Address     *url.URL
}

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
		restoreRoute(r)
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

func startApp(r *http.Request) {
	app, err := appName(r)
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("app name: %s\n", app)

	startAppUrl := fmt.Sprintf("http://%s/apps/%s/start", app, os.Getenv("TSURU_HOST"))
	authToken := fmt.Sprintf("bearer %s", os.Getenv("TOKEN"))

	client := &http.Client{}
	req, _ := http.NewRequest("POST", startAppUrl, nil)
	req.Header.Add("Authorization", authToken)
	resp, _ := client.Do(req)
	if resp.StatusCode != 200 {
		log.Printf("Error trying to start app %s", app)
		return
	}

	log.Printf("app %s started", app)
}

func getAppUnits(request *http.Request) []Unit {
	a, err := appName(request)
	if err != nil {
		log.Fatalf("Error: %s", err)
		return nil
	}
	client := &http.Client{}
	apiUrl := fmt.Sprintf("http://%s/apps/%s", a, os.Getenv("TSURU_HOST"))
	req, _ := http.NewRequest("GET", apiUrl, nil)
	authToken := fmt.Sprintf("bearer %s", os.Getenv("TOKEN"))
	req.Header.Add("Authorization", authToken)
	resp, _ := client.Do(req)

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Println("Error retrieving" + a + " api data")
		return nil
	}
	data, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatalf("Error: %s", err)
		return nil
	}

	var app App
	err1 := json.Unmarshal(data, &app)
	if err1 != nil {
		log.Fatalf("Error: %s", err)
		return nil
	}

	var units []Unit
	for _, unit := range app.Units {
		if unit.ProcessName == "web" {
			fmt.Printf("%s -> %s\n", unit.AppName, unit.Address)
			units = append(units, unit)
		}
	}
	return units
}

func restoreRoute(r *http.Request) {
	conn := redisPool.Get()
	defer conn.Close()

	name := HIPACHE_PREFIX + r.Host
	log.Print("Deleting ", name)
	_, err := conn.Do("LTRIM", name, 0, 0)
	if err != nil {
		log.Printf("Err: %s", err)
	}
}

func appName(r *http.Request) (string, error) {
	app := strings.Split(r.Host, ".")[0]
	log.Printf("appName -> %s", app)
	if app != "tsuru-caffeine-proxy" {
		return app, nil
	}

	return "", errors.New("invalid app name")
}
