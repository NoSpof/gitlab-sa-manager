package main

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

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
