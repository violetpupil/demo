package main

import (
	"context"
	"flag"
	"log"

	"demo/cancellation"

	"go.temporal.io/sdk/client"
)

func main() {
	var workflowID string
	flag.StringVar(&workflowID, "w", cancellation.WorkflowID, "workflow id")
	flag.Parse()

	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatal("unable to create client ", err)
	}
	defer c.Close()

	we, err := c.ExecuteWorkflow(context.Background(), client.StartWorkflowOptions{
		TaskQueue: cancellation.Queue,
		ID:        workflowID,
	}, cancellation.Workflow)
	if err != nil {
		log.Fatal("unable to execute workflow ", err)
	}
	log.Printf("started workflow id %s run id %s", we.GetID(), we.GetRunID())
}
