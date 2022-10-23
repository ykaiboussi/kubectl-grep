package main

import (
	"fmt"
	"log"

	text "github.com/jedib0t/go-pretty/v6/text"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/util/flowcontrol"
)

func containerState(name string) {
	cf = genericclioptions.NewConfigFlags(true)

	restConfig, err := cf.ToRESTConfig()
	if err != nil {
		log.Fatal("failed to load Config: %w", err)
	}

	restConfig.RateLimiter = flowcontrol.NewTokenBucketRateLimiter(20, 50)
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

	if len(list) != 0 {
		printPodInfo(list)
	}

	check(name, total)
}

func check(name string, total []int) {

	var found int
	for _, v := range total {
		found += v
	}

	if found == 0 {
		fmt.Printf(`Pod: %v not found  ¯\_(ツ)_/¯`+"\n", text.FgRed.Sprint(name))
	}
}
