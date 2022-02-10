package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func sendToGitlab(gitlab_id string, gitlab_scope string, gitlab_variable string, namespace string) {
	var gitlabToken = os.Getenv("GITLAB_TOKEN")
	var gitlabBaseUrl string
	var method string = "POST"
	if os.Getenv("GITLAB_BASEURL") != "" {
		gitlabBaseUrl = os.Getenv("GITLAB_BASEURL")
	} else {
		gitlabBaseUrl = "gitlab.com"
	}
	config := extractConfig(namespace)
	urldest := "https://" + gitlabBaseUrl + "/api/v4/" + gitlab_scope + "/" + gitlab_id + "/variables"
	type Payload struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}

	data := Payload{
		Key:   gitlab_variable,
		Value: config,
	}
	payloadBytes, err := json.Marshal(data)
	logIfError(err)
	body := bytes.NewReader(payloadBytes)
	if checkIfVarExist(gitlab_id, gitlab_scope, gitlab_variable) {
		method = "PUT"
		urldest = urldest + "/" + gitlab_variable

	}
	req, err := http.NewRequest(method, urldest, body)
	logIfError(err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Private-Token", gitlabToken)
	resp, err := http.DefaultClient.Do(req)
	logIfError(err)
	defer resp.Body.Close()
	log.Println("response Status:", resp.Status)
	var verbose = os.Getenv("VERBOSITY")
	if verbose == "debug" {
		log.Println("response Headers:", resp.Header)
		bodyResp, _ := ioutil.ReadAll(resp.Body)
		log.Println("response Body:", string(bodyResp))
	}

}

func checkIfVarExist(gitlab_id string, gitlab_scope string, gitlab_variable string) bool {
	var gitlabToken = os.Getenv("GITLAB_TOKEN")
	var gitlabBaseUrl string
	var status bool = false
	if os.Getenv("GITLAB_BASEURL") != "" {
		gitlabBaseUrl = os.Getenv("GITLAB_BASEURL")
	} else {
		gitlabBaseUrl = "gitlab.com"
	}
	urldest := "https://" + gitlabBaseUrl + "/api/v4/" + gitlab_scope + "/" + gitlab_id + "/variables/" + gitlab_variable
	req, err := http.NewRequest("GET", urldest, nil)
	logIfError(err)
	req.Header.Set("Private-Token", gitlabToken)
	resp, err := http.DefaultClient.Do(req)
	logIfError(err)
	var compareStatus = strings.Split(resp.Status, " ")
	log.Println(compareStatus[0])
	if compareStatus[0] == "200" {
		status = true
	}
	return status
}
