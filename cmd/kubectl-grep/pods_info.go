package main

import (
	"fmt"

	"github.com/jedib0t/go-pretty/table"
)

func printPodInfo(list []runningPods) {
	tw := table.NewWriter()

	tw.AppendHeader(table.Row{"Namespace", "Pod", "Status"})
	for _, pod := range list {
		tw.AppendRows([]table.Row{
			{pod.namespace, pod.name, pod.status},
		})
	}

	fmt.Printf("\n%s\n", tw.Render())
}
