package main

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func getSa(clientset *kubernetes.Clientset, namespace string, saName string) (string, v1.Time) {
	if saName != "" {
		saCreated, err := clientset.CoreV1().ServiceAccounts(namespace).Get(context.TODO(), saName, metav1.GetOptions{})
		if err != nil {
			panic(err)
		}
		return saCreated.GetName(), saCreated.ObjectMeta.CreationTimestamp
	}
	return saName, v1.Now()
}
