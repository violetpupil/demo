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
	flag.StringVar(&workflowID, "w", workflow.WorkflowID, "workflow id")
	flag.Parse()

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
	we, err := c.ExecuteWorkflow(context.Background(), options, workflow.Workflow, wakeUpTime)
	if err != nil {
		log.Fatal("unable to execute workflow ", err)
	}
	log.Printf("started workflow id %s run id %s, execute time %v", we.GetID(), we.GetRunID(), wakeUpTime)
}
