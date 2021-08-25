package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"path/filepath"

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
	count1 := countSa(clientset, "")
	fmt.Printf("Count of SA : %d ", count1)
	createSa(clientset, "default", "seitosan")
	count2 := countSa(clientset, "")
	fmt.Printf("Count of SA : %d ", count2)
	name, createdTime := getSa(clientset, "default", "seitosan")
	fmt.Println(name + " Created at " + createdTime.String())
	eventList, err := clientset.CoreV1().Events("").List(context.TODO(), metav1.ListOptions{
		TypeMeta: metav1.TypeMeta{
			Kind: "ServiceAccount",
		},
	})
	if err != nil {
		log.Panic(err)
	}
	fmt.Println(eventList)
	deleteSa(clientset, "default", "seitosan")
}
