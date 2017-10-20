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

	ui.HandleCurrentRowSelectionChange(func(c ecs.ContainerInstance) {
		ch := pollster.GetTaskAttributesChannel(CLUSTER_NAME, *c.ContainerInstanceArn)
		for {
			tasks := <-ch
			ui.UpdateTasks(tasks)
		}
	})

	ui.Run()
}
