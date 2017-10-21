package pollster

import (
	"fmt"
	"time"

	"github.com/onoffleftright/ecs-roll/regex"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
)

func GetContainerInstancesChannel(clusterName string) <-chan []ecs.ContainerInstance {
	c := make(chan []ecs.ContainerInstance)

	updateRate := time.Duration(1) * time.Second

	go func() {
		for {
			instanceArns, err := getContainerInstanceArns(clusterName)
			if err != nil {
				fmt.Println(err)
			}

			instanceIds := make([]string, len(instanceArns))
			for i, v := range instanceArns {
				instanceId, err := regex.ParseContainerInstanceId(v)
				if err != nil {
					fmt.Println(err)
					continue
				}

				instanceIds[i] = instanceId
			}

			containerInstances, err := getContainerInstancesAttributes(clusterName, instanceIds)
			if err != nil {
				fmt.Println(err)
				continue
			}

			c <- containerInstances
			time.Sleep(updateRate)
		}
	}()

	return c
}

func getContainerInstanceArns(clusterName string) ([]string, error) {
	results := []string{}

	if len(clusterName) == 0 {
		return results, fmt.Errorf("received zero value for argument: clusterName")
	}

	err := svc.ListContainerInstancesPages(&ecs.ListContainerInstancesInput{
		Cluster: aws.String(clusterName),
	}, func(page *ecs.ListContainerInstancesOutput, lastPage bool) bool {
		for _, v := range page.ContainerInstanceArns {
			results = append(results, *v)
		}

		return !lastPage
	})

	if err != nil {
		return results, err
	}

	return results, nil
}

func getContainerInstancesAttributes(clusterName string, containerInstanceIds []string) ([]ecs.ContainerInstance, error) {
	results := []ecs.ContainerInstance{}

	if len(clusterName) == 0 {
		return results, fmt.Errorf("received zero value for argument: clusterName")
	}

	if len(containerInstanceIds) == 0 {
		return results, fmt.Errorf("received empty slice for argument: containerInstanceIds")
	}

	response, err := svc.DescribeContainerInstances(&ecs.DescribeContainerInstancesInput{
		Cluster:            aws.String(clusterName),
		ContainerInstances: aws.StringSlice(containerInstanceIds),
	})
	if err != nil {
		return results, fmt.Errorf("could not query the describe container instances API: %v", err)
	}

	if len(response.Failures) > 0 {
		return results, fmt.Errorf("errors occured while querying the describe container instances API: %v", response.Failures)
	}

	for _, v := range response.ContainerInstances {
		results = append(results, *v)
	}

	return results, nil
}
