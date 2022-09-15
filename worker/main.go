package main

import (
	"log"

	"demo/cancellation"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatal("unable to create client ", err)
	}
	defer c.Close()

	w := worker.New(c, cancellation.Queue, worker.Options{})
	w.RegisterWorkflow(cancellation.Workflow)
	w.RegisterActivity(&cancellation.Activities{})
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatal("unable to start worker ", err)
	}
}
