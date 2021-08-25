package main

import (
	"context"

	"k8s.io/client-go/kubernetes"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func deleteSa(clientset *kubernetes.Clientset, namespace string, saName string) {

	err := clientset.CoreV1().ServiceAccounts(namespace).Delete(context.TODO(), saName, metav1.DeleteOptions{})
	if err != nil {
		panic(err)
	}
}
