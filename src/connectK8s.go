package main

import (
	"flag"
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func connectKubernetes() *kubernetes.Clientset {
	var kubeconfig *string
	if os.Getenv("RUN_IN_KATE") != "" {
		config, err := rest.InClusterConfig()
		logIfError(err)
		clientset, err := kubernetes.NewForConfig(config)
		logIfError(err)
		return clientset
	} else {
		if home := homedir.HomeDir(); home != "" {
			kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
		} else {
			kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
		}
		flag.Parse()
		// use the current context in kubeconfig
		config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
		logIfError(err)
		clientset, err := kubernetes.NewForConfig(config)
		logIfError(err)
		return clientset
	}
}
