package main

import (
	"context"
	"log"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func generateSaAuth(sa *v1.ServiceAccount, clientKate *kubernetes.Clientset) (bool, error) {
	log.Println("===> Get Service Account information about " + sa.Name)
	user, err := clientKate.CoreV1().ServiceAccounts(sa.Namespace).Get(context.TODO(), sa.Name, metav1.GetOptions{})
	logIfError(err)
	for _, ref := range user.Secrets {
		secret, err := clientKate.CoreV1().Secrets(sa.Namespace).Get(context.TODO(), ref.Name, metav1.GetOptions{})
		logIfError(err)
		server := clientKate.RESTClient().Get().URL().Host
		if server != "" {
			generateKubeConfig(secret, sa.Namespace, server, sa.Name)
		}
	}

	return false, nil
}

func generateKubeConfig(secret *v1.Secret, namespace string, server string, saName string) {
	clusters := make(map[string]*clientcmdapi.Cluster)
	clusters["default-cluster"] = &clientcmdapi.Cluster{
		Server:                   "https://" + server,
		CertificateAuthorityData: secret.Data["ca.crt"],
	}

	contexts := make(map[string]*clientcmdapi.Context)
	contexts["default-context"] = &clientcmdapi.Context{
		Cluster:   "default-cluster",
		Namespace: namespace,
		AuthInfo:  saName,
	}

	authinfos := make(map[string]*clientcmdapi.AuthInfo)
	authinfos[saName] = &clientcmdapi.AuthInfo{
		Token: string(secret.Data["token"]),
	}

	clientConfig := clientcmdapi.Config{
		Kind:           "Config",
		APIVersion:     "v1",
		Clusters:       clusters,
		Contexts:       contexts,
		CurrentContext: "default-context",
		AuthInfos:      authinfos,
	}
	clientcmd.WriteToFile(clientConfig, "./"+namespace+".kubeconfig")
}
