package main

import (
	"log"

	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/dynamic"
)

func main() {
	cf = genericclioptions.NewConfigFlags(true)

	restConfig, err := cf.ToRESTConfig()
	if err != nil {
		log.Fatal("failed to load Config: %w", err)
	}

	dyn, err := dynamic.NewForConfig(restConfig)
	if err != nil {
		log.Fatal("failed to construct dynamic client: %w", err)
	}

	dc, err := cf.ToDiscoveryClient()
	if err != nil {
		log.Fatal("failed to discover client: %w", err)
	}

	apis, err := findAPIs(dc)
	if err != nil {
		log.Fatal("failed to load APIs: %w", err)
	}

	l, err := getAllResources(dyn, apis.resources(), true)
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range l {
		lookupPod("9492j", v.Object)
	}
}
