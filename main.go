package main

import (
	"github.com/onoffleftright/ecs-roll/pollster"
	"github.com/onoffleftright/ecs-roll/ui"
)

func main() {
	go func() {
		c := pollster.GetContainerInstancesChannel("microservices_20170925")
		for {
			containerInstances := <-c
			ui.UpdateContainerInstances(containerInstances)
		}
	}()

	ui.Run()
}
