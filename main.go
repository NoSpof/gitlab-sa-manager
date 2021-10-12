package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main() {
	clientKate := connectKubernetes()
	//	var gitlab = os.Getenv("GITLAB_TOKEN")
	//	var env = os.Getenv("ENV")

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
		gitlab_id, gitlab_scope, gitlab_variable, namespace := getInfo(annotation, sa, event)
		fmt.Println(gitlab_id)
		fmt.Println(gitlab_scope)
		fmt.Println(gitlab_variable)
		if event.Type == "ADDED" {
			log.Println("===> Get Service Account information about " + sa.Name)
			sa_secrets, err := clientKate.CoreV1().ServiceAccounts(namespace).Get(context.TODO(), sa.Name, metav1.GetOptions{})
			if err != nil {
				panic(err)
			}
			sa_token, err := clientKate.CoreV1().Secrets(namespace).Get(context.TODO(), sa_secrets.Secrets[0].Name, metav1.GetOptions{})
			if err != nil {
				panic(err)
			}
			token := sa_token.Data["token"]
			//ca_crt := sa_token.Data["ca.crt"]
			tokenDec, err := base64.StdEncoding.DecodeString(string(token))
			if err != nil {
				panic(err)
			}
			fmt.Println(string(tokenDec))
		}
	}
}
