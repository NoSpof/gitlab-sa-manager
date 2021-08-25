package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"path/filepath"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	sa, err := clientset.CoreV1().ServiceAccounts("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%d servives acclunts founded ", len(sa.Items))

	createdSA := &v1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "seitosan",
			Namespace: "default",
		},
	}
	result, err := clientset.CoreV1().ServiceAccounts("default").Create(context.TODO(), createdSA, metav1.CreateOptions{})
	if err != nil {
		log.Println(err)
	}
	fmt.Printf("Created sa %q.\n", result.GetObjectMeta().GetName())

	sa2, err := clientset.CoreV1().ServiceAccounts("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%d servives accounts founded ", len(sa2.Items))

	saCreated, err := clientset.CoreV1().ServiceAccounts("default").Get(context.TODO(), "seitosan", metav1.GetOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Println(saCreated.Secrets)
	err = clientset.CoreV1().ServiceAccounts("default").Delete(context.TODO(), "seitosan", metav1.DeleteOptions{})
	if err != nil {
		panic(err)
	}

}
