package main

import "k8s.io/cli-runtime/pkg/genericclioptions"

var cf *genericclioptions.ConfigFlags

func getNamespace() string {
	if v := *cf.Namespace; v != "" {
		return v
	}
	clientConfig := cf.ToRawKubeConfigLoader()
	defaultNamespace, _, err := clientConfig.Namespace()
	if err != nil {
		defaultNamespace = "default"
	}
	return defaultNamespace
}
