package pollster

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
)

var svc *ecs.ECS

func init() {
	svc = ecs.New(session.New())
}

func GetContainerInstancesChannel(clusterName string) <-chan []string {
	c := make(chan []string)

	go func() {
		for {
			instances, err := GetContainerInstances(clusterName)
			if err != nil {
				fmt.Println(err)
			}

			c <- instances
		}
	}()

	return c
}

func GetContainerInstances(clusterName string) ([]string, error) {
	results := []string{}

	if len(clusterName) == 0 {
		return results, fmt.Errorf("received zero value for clusterName param")
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
