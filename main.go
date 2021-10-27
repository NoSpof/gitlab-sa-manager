package main

import (
	"context"
	"fmt"
	"log"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main() {
	clientKate := connectKubernetes()

	sa, err := clientKate.CoreV1().ServiceAccounts("").Watch(context.TODO(), metav1.ListOptions{LabelSelector: "gitlab.sa-manager.k8s.io/enable=true"})
	if err != nil {
		panic(err)
	}
	for event := range sa.ResultChan() {

		sa, ok := event.Object.(*v1.ServiceAccount)
		if !ok {
			log.Println("Unexpected Type")
		}
		annotation := sa.ObjectMeta.GetAnnotations()
		if annotation["sa-manager.k8s.io"] != "" {
			switch annotation["sa-manager.k8s.io"] {
			case "gitlab":
				fmt.Println("Repository destination gitlab")
				gitlab_id, gitlab_scope, gitlab_variable, namespace := getInfoFromGitlab(annotation, sa, event)
				log.Println(gitlab_id)
				log.Println(gitlab_scope)
				log.Println(gitlab_variable)
				log.Println(namespace)
			case "github":
				fmt.Println("Repository destination github")
			case "azDevops":
				fmt.Println("Repository destination Azure Devops")
			default:
				log.Println("Annotation sa-manager.k8s.io must be set with : gitlab github or azDevops parameter")
			}
		}
		fmt.Println(event.Type)
		switch event.Type {
		case "ADDED":
			generateSaAuth(sa, clientKate)
		case "DELETED":
			log.Println("Deleting sa token in repo")
		case "MODIFIED":

		default:
			log.Println(event.Type + " not supported at time ")
		}

	}
}
