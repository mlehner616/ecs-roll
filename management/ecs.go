package management

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
)

var svc *ecs.ECS

func init() {
	svc = ecs.New(session.New())
}

type State string

const (
	Active   State = "ACTIVE"
	Draining State = "DRAINING"
)

func SetContainerInstanceState(clusterName, containerInstanceArn string, s State) error {
	if len(clusterName) == 0 {
		return errors.New("received zero value for argument: clusterName")
	}

	if len(containerInstanceArn) == 0 {
		return errors.New("received zero value for argument: containerInstanceArn")
	}

	results, err := svc.UpdateContainerInstancesState(&ecs.UpdateContainerInstancesStateInput{
		Cluster:            aws.String(clusterName),
		ContainerInstances: aws.StringSlice([]string{containerInstanceArn}),
		Status:             aws.String(string(s)),
	})
	if err != nil {
		return fmt.Errorf("could not query the update container instance state API: %v", err)
	}

	if len(results.Failures) != 0 {
		return fmt.Errorf("an error occurred while attempting to update the container instance state: %v", results.Failures)
	}

	return nil
}
