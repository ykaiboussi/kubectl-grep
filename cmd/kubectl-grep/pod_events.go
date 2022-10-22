package main

import (
	"context"
	"fmt"
	"log"

	table "github.com/jedib0t/go-pretty/table"
	text "github.com/jedib0t/go-pretty/v6/text"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubernetes "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func getPodEvents(restConfig *rest.Config, namespace, podName, status string) {
	clientset, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		log.Fatal(err)
	}
	events, err := clientset.CoreV1().Events(namespace).List(context.TODO(), metav1.ListOptions{FieldSelector: "involvedObject.name=" + podName, TypeMeta: metav1.TypeMeta{Kind: "Pod"}})
	if err != nil {
		log.Fatal(err)
	}
	tw := table.NewWriter()
	// append a header row
	tw.SetTitle("Status: %v %v", podName, text.FgRed.Sprint(status))
	tw.AppendHeader(table.Row{"Time", "Reason", "Message"})
	for _, item := range events.Items {
		// append some data rows
		tw.AppendRows([]table.Row{
			{item.LastTimestamp, item.Reason, text.WrapSoft(item.Message, 50)},
		})
	}
	fmt.Printf("Output:\n%s\n", tw.Render())
}
