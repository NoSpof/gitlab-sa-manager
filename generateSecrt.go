package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/watch"
)

type GithubBody struct {
	encrypted_value string
}

func generateSecret(annotation map[string]string, sa *v1.ServiceAccount, event watch.Event) {
	if annotation["sa-manager.k8s.io"] != "" {
		switch annotation["sa-manager.k8s.io"] {
		case "gitlab":
			fmt.Println("Repository destination gitlab")
			gitlab_id := annotation["gitlab.sa-manager.k8s.io/id"]
			gitlab_scope := annotation["gitlab.sa-manager.k8s.io/scope"]
			gitlab_variable := annotation["gitlab.sa-manager.k8s.io/variable"]
			namespace := sa.Namespace
			var verbose = os.Getenv("VERBOSITY")
			if verbose == "debug" {
				log.Printf("Service Account %s  has been  %v on namespace %s\n", sa.Name, event.Type, sa.Namespace)
				log.Println("ID Gitlab :" + gitlab_id)
				log.Println("Scope Gitlab :" + gitlab_scope)
				log.Println("Variable Gitlab : " + gitlab_variable)
				log.Println("Namespace : " + namespace)
			}
			sendToGitlab(gitlab_id, gitlab_scope, gitlab_variable, namespace)
		case "github":
			//var github_baseUrl = os.Getenv("GITHUB_BASEURL")
			var github_baseUrl = "https://api.github.com/"
			fmt.Println("Repository destination github")
			github_owner := annotation["github.sa-manager.k8s.io/owner"]
			github_repo := annotation["github.sa-manager.k8s.io/repository"]
			github_variable := annotation["github.sa-manager.k8s.io/variable"]
			github_scope := annotation["github.sa-manager.k8s.io/scope"]
			namespace := sa.Namespace
			var verbose = os.Getenv("VERBOSITY")
			if verbose == "debug" {
				log.Printf("Service Account %s  has been  %v on namespace %s\n", sa.Name, event.Type, sa.Namespace)
				log.Println("Owner Github :" + github_owner)
				log.Println("Scope Gitlab :" + github_scope)
				log.Println("Variable Github: " + github_variable)
				log.Println("Repository considerated : " + github_repo)
				log.Println("Namespace : " + namespace)
			}
			if github_scope == "repo" {
				client := &http.Client{}
				data := GithubBody{
					encrypted_value: "bdbchchcbccbhc==",
				}
				json, err := json.Marshal(data)
				if err != nil {
					panic(err)
				}

				// set the HTTP method, url, and request body
				req, err := http.NewRequest(http.MethodPut, github_baseUrl+"/repos/"+github_owner+"/"+github_repo+"/actions/secrets/"+github_variable, bytes.NewBuffer(json))
				if err != nil {
					panic(err)
				}

				// set the request header Content-Type for json
				req.Header.Set("Accept", "application/vnd.github.v3+json")
				resp, err := client.Do(req)
				if err != nil {
					panic(err)
				}

				fmt.Println(resp.StatusCode)
			}
		case "azDevops":
			fmt.Println("Repository destination Azure Devops")
		default:
			log.Println("Annotation sa-manager.k8s.io must be set with : gitlab github or azDevops parameter")
		}
	}

}
