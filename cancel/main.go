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

	err = c.CancelWorkflow(context.Background(), workflowID, "")
	if err != nil {
		log.Fatal("unable to cancel workflow ", err)
	}
	log.Print("cancel workflow id ", cancellation.WorkflowID)
}
