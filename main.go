package main

import (
	"github.com/aws/aws-sdk-go/service/ecs"

	"github.com/onoffleftright/ecs-roll/pollster"
	"github.com/onoffleftright/ecs-roll/ui"
)

const CLUSTER_NAME = "microservices_20170925"

func main() {
	go func() {
		c := pollster.GetContainerInstancesChannel(CLUSTER_NAME)
		for {
			containerInstances := <-c
			ui.UpdateContainerInstances(containerInstances)
		}
	}()

	var activeTaskAttributesDataChan <-chan []ecs.Task
	var activeTaskAttributesStopChan chan<- bool
	ui.HandleCurrentRowSelectionChange(func(c ecs.ContainerInstance) {
		if activeTaskAttributesStopChan != nil {
			activeTaskAttributesStopChan <- true
		}

		activeTaskAttributesDataChan, activeTaskAttributesStopChan = pollster.GetTaskAttributesChannel(CLUSTER_NAME, *c.ContainerInstanceArn)
		for {
			tasks := <-activeTaskAttributesDataChan
			ui.UpdateTasks(tasks)
		}
	})

	ui.Run()
}
