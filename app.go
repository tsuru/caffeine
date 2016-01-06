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
	if app == CAFFEINE_APP_NAME {
		return "", errors.New("invalid app name")
	}

	appName, err := getAppName(app)
	if err == nil {
		return appName, nil
	}

	appName, err = getAppNameByCname(app)
	if err == nil {
		return appName, nil
	}

	return "", err
}

func getAppName(appName string) (string, error) {
	listAppsUrl := fmt.Sprintf("%s/apps/?name=%s", os.Getenv("TSURU_HOST"), appName)
	authToken := fmt.Sprintf("bearer %s", os.Getenv("TOKEN"))

	client := &http.Client{}
	req, _ := http.NewRequest("GET", listAppsUrl, nil)
	req.Header.Add("Authorization", authToken)
	resp, _ := client.Do(req)
	if resp.StatusCode != 200 {
		log.Printf("Error trying to get app %s", appName)
		return "", errors.New("Error trying to get app info")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error trying to get app %s", appName)
		return "", errors.New("Error trying to get app info")
	}

	var apps []App
	json.Unmarshal(body, &apps)
	if len(apps) == 0 {
		log.Printf("App %s not found", appName)
		return "", errors.New("App not found")
	}

	return apps[0].Name, nil
}

func getAppNameByCname(appName string) (string, error) {
	listAppsUrl := fmt.Sprintf("%s/apps/", os.Getenv("TSURU_HOST"))
	authToken := fmt.Sprintf("bearer %s", os.Getenv("TOKEN"))

	client := &http.Client{}
	req, _ := http.NewRequest("GET", listAppsUrl, nil)
	req.Header.Add("Authorization", authToken)
	resp, _ := client.Do(req)
	if resp.StatusCode != 200 {
		log.Printf("Error trying to get app %s", appName)
		return "", errors.New("Error trying to get app info")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error trying to get app %s", appName)
		return "", errors.New("Error trying to get app info")
	}

	var apps []App
	err = json.Unmarshal(body, &apps)
	if err != nil {
		log.Printf("Error trying to get app %s", appName)
		return "", errors.New("Error trying to get app info")
	}
	for _, app := range apps {
		for _, cname := range app.Cname {
			if cname == appName {
				return app.Name, nil
			}
		}
	}

	log.Printf("App %s not found", appName)
	return "", errors.New("App not found")
}
