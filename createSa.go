package main

import (
	"context"
	"fmt"
	"log"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func createSa(clientset *kubernetes.Clientset, namespace string, saName string) {

	createdSA := &v1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      saName,
			Namespace: namespace,
		},
	}
	result, err := clientset.CoreV1().ServiceAccounts(namespace).Create(context.TODO(), createdSA, metav1.CreateOptions{})
	if err != nil {
		log.Println(err)
	}
	fmt.Printf("Created sa %q.\n", result.GetObjectMeta().GetName())

}
