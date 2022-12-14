package main

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	text "github.com/jedib0t/go-pretty/v6/text"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/klog"
)

var (
	// adding a counter goroutine to check if pod is present
	i     = 0
	total []int

	// adding a goroutine to append Running Pods information
	p    runningPods
	list []runningPods

	wg sync.WaitGroup
)

func getAllResources(client dynamic.Interface, apis []apiResource, allNs bool) ([]unstructured.Unstructured, error) {
	var mu sync.Mutex
	var wg sync.WaitGroup
	var out []unstructured.Unstructured

	start := time.Now()
	klog.V(2).Infof("starting to query %d APIs in concurrently", len(apis))

	var errResult error
	for _, api := range apis {
		if !allNs && !api.r.Namespaced {
			klog.V(4).Infof("[query api] api (%s) is non-namespaced, skipping", api.r.Name)
			continue
		}
		wg.Add(1)
		go func(a apiResource) {
			defer wg.Done()
			klog.V(4).Infof("[query api] start: %s", a.GroupVersionResource())
			v, err := queryAPI(client, a, allNs)
			if err != nil {
				klog.V(4).Infof("[query api] error querying: %s, error=%v", a.GroupVersionResource(), err)
				errResult = err
				return
			}
			mu.Lock()
			out = append(out, v...)
			mu.Unlock()
			klog.V(4).Infof("[query api]  done: %s, found %d apis", a.GroupVersionResource(), len(v))
		}(api)
	}

	klog.V(2).Infof("fired up all goroutines to query APIs")
	wg.Wait()
	klog.V(2).Infof("all goroutines have returned in %v", time.Since(start))
	klog.V(2).Infof("query result: error=%v, objects=%d", errResult, len(out))
	return out, errResult
}

func queryAPI(client dynamic.Interface, api apiResource, allNs bool) ([]unstructured.Unstructured, error) {
	var out []unstructured.Unstructured

	var next string
	var ns string

	if !allNs {
		ns = getNamespace()
	}
	for {
		var intf dynamic.ResourceInterface
		nintf := client.Resource(api.GroupVersionResource())
		if !allNs {
			intf = nintf.Namespace(ns)
		} else {
			intf = nintf
		}
		resp, err := intf.List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return nil, fmt.Errorf("listing resources failed (%s): %w", api.GroupVersionResource(), err)
		}
		out = append(out, resp.Items...)

		next = resp.GetContinue()
		if next == "" {
			break
		}
	}
	return out, nil
}

func getPodName(k map[string]interface{}) string {
	m, ok := k["metadata"].(map[string]interface{})
	if !ok {
		return ""
	}
	return m["name"].(string)
}

func lookupNamespace(k map[string]interface{}) string {
	m, ok := k["metadata"].(map[string]interface{})
	if !ok {
		return "no namespace assign to the pod"
	}
	return m["namespace"].(string)
}

func increment() {
	i += 1
	total = append(total, i)
}

func appendRunningPodds(podName, namespace string) {
	defer wg.Done()
	p.name = podName
	p.namespace = namespace
	p.status = text.FgGreen.Sprint("Running")
	list = append(list, p)
}
func lookupPod(restConfig *rest.Config, k map[string]interface{}, input string) {
	podName := getPodName(k)
	if strings.Contains(podName, input) {
		namespace := lookupNamespace(k)
		containerInfo, _ := k["status"].(map[string]interface{})
		containerStatuses, _ := containerInfo["containerStatuses"].([]interface{})
		for _, c := range containerStatuses {
			containerStateObj, _ := c.(map[string]interface{})
			state, _ := containerStateObj["state"].(map[string]interface{})
			_, ok := state["running"]
			if !ok {
				go increment()
				getPodEvents(namespace, podName, "Not running")
			} else {
				go increment()

				wg.Add(1)
				go appendRunningPodds(podName, namespace)
			}
		}
	}
	wg.Wait()
}
