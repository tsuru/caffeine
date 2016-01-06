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
	Ip    string
	Cname []string
}

func startApp(hostname string) {
	app, err := appName(hostname)
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

func appName(hostname string) (string, error) {
	appName := strings.Split(hostname, ".")[0]
	if appName == CAFFEINE_APP_NAME {
		return "", errors.New("invalid app name")
	}

	app, err := getApp(hostname)
	if err == nil {
		return app.Name, nil
	}

	return "", err
}

func getApp(hostname string) (*App, error) {
	apps, err := listApps(nil)
	if err != nil {
		log.Println("Error trying to get apps info")
		return nil, err
	}

	for _, app := range apps {
		if hostname == app.Ip {
			return &app, nil
		}
		for _, cname := range app.Cname {
			if hostname == cname {
				return &app, nil
			}
		}
	}

	log.Printf("App %s not found", hostname)
	return nil, errors.New("App not found")
}

func authToken() string {
	return fmt.Sprintf("bearer %s", os.Getenv("TOKEN"))
}

func listApps(queryParams map[string]string) ([]App, error) {
	listAppsURL := fmt.Sprintf("%s/apps/%s", os.Getenv("TSURU_HOST"), queryParamsToString(queryParams))

	client := &http.Client{}
	req, _ := http.NewRequest("GET", listAppsURL, nil)
	req.Header.Add("Authorization", authToken())
	resp, _ := client.Do(req)
	if resp.StatusCode != 200 {
		return nil, errors.New("Error trying to get app info")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("Error trying to get app info")
	}

	var apps []App
	err = json.Unmarshal(body, &apps)
	if err != nil {
		return nil, errors.New("Error trying to get app info")
	}

	return apps, nil
}

func queryParamsToString(queryParams map[string]string) string {
	str := ""
	for key, value := range queryParams {
		var separator string
		if len(str) == 0 {
			separator = "?"
		} else {
			separator = "&"
		}

		str = fmt.Sprintf("%s%s%s=%s", str, separator, key, value)
	}
	return str
}
