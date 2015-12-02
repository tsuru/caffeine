package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

func startApp(r *http.Request) {
	app, err := appName(r)
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

func appName(r *http.Request) (string, error) {
	app := strings.Split(r.Host, ".")[0]
	log.Printf("appName -> %s", app)
	if app != "tsuru-caffeine-proxy" {
		return app, nil
	}

	return "", errors.New("invalid app name")
}
