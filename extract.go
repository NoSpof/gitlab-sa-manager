package main

import (
	b64 "encoding/base64"
	"io/ioutil"
)

func extractConfig(namespace string) string {
	content, err := ioutil.ReadFile("./" + namespace + ".kubeconfig")
	logIfError(err)
	text := string(content)
	sEnc := b64.StdEncoding.EncodeToString([]byte(text))
	return sEnc
}
