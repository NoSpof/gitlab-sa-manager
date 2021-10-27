package main

import (
	"context"
	"log"
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func generateSaAuth(sa *v1.ServiceAccount, clientKate *kubernetes.Clientset) {
	log.Println("===> Get Service Account information about " + sa.Name)
	saName, err := clientKate.CoreV1().ServiceAccounts(sa.Namespace).Get(context.TODO(), sa.Name, metav1.GetOptions{})
	logIfError(err)
	// waiting for sa creation
	time.Sleep(10 * time.Second)
	saSecret, err := clientKate.CoreV1().Secrets(sa.Namespace).Get(context.TODO(), saName.Secrets[0].Name, metav1.GetOptions{})
	logIfError(err)
	server := clientKate.RESTClient().Get().URL().Host
	if server != "" {
		generateKubeConfig(saSecret, sa.Namespace, server, sa.Name)
	}
}
