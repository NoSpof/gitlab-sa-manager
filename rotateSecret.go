package main

import (
	"context"
	"log"
	"time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func rotateSecret(clientset *kubernetes.Clientset, namespace string, saName string) {

	name, tokenName, createdTime := getSa(clientset, "seitosan", "seitosan")
	log.Println(name + " Created at " + createdTime.String() + " Token : " + tokenName)
	secret, err := clientset.CoreV1().Secrets(namespace).Get(context.TODO(), tokenName, v1.GetOptions{})
	if err != nil {
		log.Panic(err)
	}
	secretCreationDate := secret.ObjectMeta.CreationTimestamp
	today := time.Now()
	insec := secretCreationDate.Add(addingTime * time.Minute)
	//Check Expiration date
	log.Println(saName + " Expired at : " + insec.String())
	if today.After(insec) {
		log.Println("Service account token : " + tokenName + " expired")
		log.Println(insec.String() + " After " + today.String())
		err := clientset.CoreV1().Secrets(namespace).Delete(context.TODO(), tokenName, v1.DeleteOptions{})
		if err != nil {
			log.Panic(err)
		}
		log.Println("Rotate for " + tokenName + " Success")
	} else {
		log.Println("Service account token for " + tokenName + " not expired")
	}

}
