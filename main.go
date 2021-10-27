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
		switch annotation["sa-manager.k8s.io"] {
		case "gitlab":
			fmt.Println("Repository destination gitlab")
		case "github":
			fmt.Println("Repository destination github")
		case "azDevops":
			fmt.Println("Repository destination Azure Devops")
		default:
			log.Println("Annotation sa-manager.k8s.io must be set with : gitlab github or azDevops parameter")
		}
		if annotation["sa-manager.k8s.io"] != "" {
			gitlab_id, gitlab_scope, gitlab_variable, namespace := getInfoFromGitlab(annotation, sa, event)
			log.Println(gitlab_id)
			log.Println(gitlab_scope)
			log.Println(gitlab_variable)

			if event.Type == "ADDED" {
				log.Println("===> Get Service Account information about " + sa.Name)
				sa_secrets, err := clientKate.CoreV1().ServiceAccounts(namespace).Get(context.TODO(), sa.Name, metav1.GetOptions{})
				logIfError(err)
				sa_token, err := clientKate.CoreV1().Secrets(namespace).Get(context.TODO(), sa_secrets.Secrets[0].Name, metav1.GetOptions{})
				logIfError(err)
				data := sa_token.Data
				ca_crt := data["ca.crt"]
				token := data["token"]
				


				fmt.Println(string(ca_crt))
				fmt.Println(string(token))

			}

		}
	}
}
