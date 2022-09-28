package main

import (
	"context"
	"demo/workflow"
	"log"
	"time"

	"go.temporal.io/sdk/client"
)

func main() {
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatal("unable to create client ", err)
	}
	defer c.Close()

	wakeUpTime := time.Now().Add(20 * time.Second)
	err = c.SignalWorkflow(context.Background(), workflow.WorkflowID, "", workflow.SignalType, wakeUpTime)
	if err != nil {
		log.Fatal("unable to signal workflow ", err)
	}
	log.Print("signal workflow ", workflow.WorkflowID, " wake up time ", wakeUpTime)
}
