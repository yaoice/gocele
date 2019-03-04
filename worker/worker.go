package main

import (
	"github.com/yaoice/gocele/pkg/controller"
	"github.com/yaoice/gocele/sample"
	"log"
)

func main() {

	calC := controller.NewCalController()
	// Register tasks
	tasks := map[string]interface{}{
		"add":        exampletasks.Add,
		"multiply":   exampletasks.Multiply,
		"panic_task": exampletasks.PanicTask,
	}
	calC.RegisterTasks(tasks)

	worker := calC.CreateWorker("machinery_worker")
	if err := worker.Launch(); err != nil {
		log.Fatal("Could not launch worker", err)
	}
}
