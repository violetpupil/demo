package main

import (
	"context"
	"demo/workflow"
	"flag"
	"log"
	"time"

	"go.temporal.io/sdk/client"
)

func main() {
	var workflowID string
	var disableUpdate bool
	flag.StringVar(&workflowID, "w", workflow.WorkflowID, "workflow id")
	flag.BoolVar(&disableUpdate, "d", false, "disable wake up time update")
	flag.Parse()
	var flow func(workflow.Context, time.Time) error
	if disableUpdate {
		flow = workflow.DelayWorkflow
	} else {
		flow = workflow.Workflow
	}

	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatal("unable to create client ", err)
	}
	defer c.Close()

	options := client.StartWorkflowOptions{
		TaskQueue: workflow.Queue,
		ID:        workflowID,
	}
	wakeUpTime := time.Now().Add(30 * time.Second)
	we, err := c.ExecuteWorkflow(context.Background(), options, flow, wakeUpTime)
	if err != nil {
		log.Fatal("unable to execute workflow ", err)
	}
	log.Printf("started workflow id %s run id %s, execute time %v", we.GetID(), we.GetRunID(), wakeUpTime)
}
