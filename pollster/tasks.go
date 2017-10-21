package pollster

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
)

func GetTaskAttributesChannel(clusterName, containerInstanceArn string) <-chan []ecs.Task {
	c := make(chan []ecs.Task)

	updateRate := time.Duration(1) * time.Second

	go func() {
		for {
			taskArns, err := getTaskArns(clusterName, containerInstanceArn)
			if err != nil {
				fmt.Println(err)
				continue
			}

			// Nothing to do this loop
			if len(taskArns) == 0 {
				continue
			}

			tasks, err := getTaskAttributes(clusterName, taskArns)
			if err != nil {
				fmt.Println(err)
				continue
			}

			c <- tasks
			time.Sleep(updateRate)
		}
	}()

	return c
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
