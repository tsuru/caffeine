package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

const CAFFEINE_APP_NAME = "tsuru-caffeine-proxy"

func startApp(host string) {
	app, err := appName(host)
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("app name: %s\n", app)

	startAppUrl := fmt.Sprintf("%s/apps/%s/start", os.Getenv("TSURU_HOST"), app)
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

func appName(host string) (string, error) {
	app := strings.Split(host, ".")[0]
	if app != CAFFEINE_APP_NAME {
		return app, nil
	}

	return "", errors.New("invalid app name")
}
