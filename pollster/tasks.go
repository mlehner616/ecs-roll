package pollster

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
)

func GetTaskAttributesChannel(clusterName, containerInstanceArn string) (<-chan []ecs.Task, chan<- bool) {
	data := make(chan []ecs.Task)
	stop := make(chan bool)

	updateRate := time.Duration(1) * time.Second
	ticker := time.Tick(updateRate)

	go func() {
		// Send immediately for the initial load, don't wait on ticker
		data <- getTasksAndHandleErrors(clusterName, containerInstanceArn)

		for {
			select {
			case <-stop:
				return
			case <-ticker:
				data <- getTasksAndHandleErrors(clusterName, containerInstanceArn)
			}
		}
	}()

	return data, stop
}

func getTasksAndHandleErrors(clusterName, containerInstanceArn string) []ecs.Task {
	taskArns, err := getTaskArns(clusterName, containerInstanceArn)
	if err != nil {
		fmt.Println(err)
		return []ecs.Task{}
	}

	// Nothing to do this loop
	if len(taskArns) == 0 {
		return []ecs.Task{}
	}

	tasks, err := getTaskAttributes(clusterName, taskArns)
	if err != nil {
		fmt.Println(err)
		return []ecs.Task{}
	}

	return tasks
}

func getTaskArns(clusterName, containerInstanceArn string) ([]string, error) {
	results := []string{}

	if len(clusterName) == 0 {
		return results, fmt.Errorf("received zero value for argument: clusterName")
	}

	if len(containerInstanceArn) == 0 {
		return results, fmt.Errorf("received zero value for argument: containerInstanceArn")
	}

	response, err := svc.ListTasks(&ecs.ListTasksInput{
		Cluster:           aws.String(clusterName),
		ContainerInstance: aws.String(containerInstanceArn),
	})
	if err != nil {
		return results, fmt.Errorf("could not query the list tasks API: %v", err)
	}

	for _, v := range response.TaskArns {
		results = append(results, *v)
	}

	return results, nil
}

func getTaskAttributes(clusterName string, taskArns []string) ([]ecs.Task, error) {
	results := []ecs.Task{}

	if len(clusterName) == 0 {
		return results, fmt.Errorf("received zero value for argument: clusterName")
	}

	if len(taskArns) == 0 {
		return results, fmt.Errorf("received empty slice for argument: taskArns")
	}

	response, err := svc.DescribeTasks(&ecs.DescribeTasksInput{
		Cluster: aws.String(clusterName),
		Tasks:   aws.StringSlice(taskArns),
	})
	if err != nil {
		return results, fmt.Errorf("could not query the describe tasks API: %v", err)
	}

	for _, v := range response.Tasks {
		results = append(results, *v)
	}

	return results, nil
}
