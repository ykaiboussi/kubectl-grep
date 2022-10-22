package main

import (
	"fmt"
	"log"

	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/dynamic"
)

func containerState(name string) {
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
		lookupPod(restConfig, v.Object, name)
	}

	var found int
	for _, v := range total {
		found += v
	}

	if found == 0 {
		fmt.Printf("Pod: %v not found\n", name)
	}
}
