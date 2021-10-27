package main

import (
	"log"
	"os"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/watch"
)

func getInfoFromGitlab(annotation map[string]string, sa *v1.ServiceAccount, event watch.Event) (string, string, string, string) {
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
	return gitlab_id, gitlab_scope, gitlab_variable, namespace
}
