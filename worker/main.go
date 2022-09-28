package main

import (
	"demo/workflow"
	"log"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatal("unable to create client ", err)
	}
	defer c.Close()

	w := worker.New(c, workflow.Queue, worker.Options{})
	w.RegisterWorkflow(workflow.Workflow)
	w.RegisterWorkflow(workflow.DelayWorkflow)
	w.RegisterActivity(&workflow.Activities{})
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatal("unable to start worker ", err)
	}
}
