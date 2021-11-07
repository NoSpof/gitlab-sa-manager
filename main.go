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
		fmt.Println(event.Type)
		switch event.Type {
		case "ADDED":
			generateSaAuth(sa, clientKate)
			generateSecret(annotation, sa, event)
		case "DELETED":
			log.Println("Deleting sa token in repo")
		case "MODIFIED":
			generateSaAuth(sa, clientKate)
			generateSecret(annotation, sa, event)
		default:
			log.Println(event.Type + " not supported at time ")
		}

	}
}
