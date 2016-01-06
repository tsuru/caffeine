package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

const CAFFEINE_APP_NAME = "tsuru-caffeine-proxy"

type App struct {
	Name  string
	Cname []string
}

func startApp(host string) {
	app, err := appName(host)
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("app name: %s\n", app)

	startAppURL := fmt.Sprintf("%s/apps/%s/start", os.Getenv("TSURU_HOST"), app)

	client := &http.Client{}
	req, _ := http.NewRequest("POST", startAppURL, nil)
	req.Header.Add("Authorization", authToken())
	resp, _ := client.Do(req)
	if resp.StatusCode != 200 {
		log.Printf("Error trying to start app %s", app)
		return
	}

	log.Printf("app %s started", app)
}

func appName(host string) (string, error) {
	appName := strings.Split(host, ".")[0]
	if appName == CAFFEINE_APP_NAME {
		return "", errors.New("invalid app name")
	}

	app, err := getAppByName(appName)
	if err == nil {
		return app.Name, nil
	}

	app, err = getAppByCname(appName)
	if err == nil {
		return app.Name, nil
	}

	return "", err
}

func getAppByName(appName string) (*App, error) {
	listAppsURL := fmt.Sprintf("%s/apps/?name=%s", os.Getenv("TSURU_HOST"), appName)

	client := &http.Client{}
	req, _ := http.NewRequest("GET", listAppsURL, nil)
	req.Header.Add("Authorization", authToken())
	resp, _ := client.Do(req)
	if resp.StatusCode != 200 {
		log.Printf("Error trying to get app %s", appName)
		return nil, errors.New("Error trying to get app info")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error trying to get app %s", appName)
		return nil, errors.New("Error trying to get app info")
	}

	var apps []App
	json.Unmarshal(body, &apps)
	if len(apps) == 0 {
		log.Printf("App %s not found", appName)
		return nil, errors.New("App not found")
	}

	return &apps[0], nil
}

func getAppByCname(appName string) (*App, error) {
	listAppsURL := fmt.Sprintf("%s/apps/", os.Getenv("TSURU_HOST"))

	client := &http.Client{}
	req, _ := http.NewRequest("GET", listAppsURL, nil)
	req.Header.Add("Authorization", authToken())
	resp, _ := client.Do(req)
	if resp.StatusCode != 200 {
		log.Printf("Error trying to get app %s", appName)
		return nil, errors.New("Error trying to get app info")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error trying to get app %s", appName)
		return nil, errors.New("Error trying to get app info")
	}

	var apps []App
	err = json.Unmarshal(body, &apps)
	if err != nil {
		log.Printf("Error trying to get app %s", appName)
		return nil, errors.New("Error trying to get app info")
	}
	for _, app := range apps {
		for _, cname := range app.Cname {
			if cname == appName {
				return &app, nil
			}
		}
	}

	log.Printf("App %s not found", appName)
	return nil, errors.New("App not found")
}

func authToken() string {
	return fmt.Sprintf("bearer %s", os.Getenv("TOKEN"))
}
