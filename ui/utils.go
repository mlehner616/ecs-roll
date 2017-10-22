package ui

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/ecs"

	"github.com/onoffleftright/ecs-roll/regex"
)

func parseContainerInstanceId(containterInstanceArn string) string {
	out, err := regex.ParseContainerInstanceId(containterInstanceArn)
	if err != nil {
		fmt.Println(err)
		return "ERROR"
	}

	return out
}

func parseTaskArn(taskArn string) string {
	out, err := regex.ParseTaskArn(taskArn)
	if err != nil {
		fmt.Println(err)
		return "ERROR"
	}

	return out
}

func parseTaskDefinitionArn(taskDefinitionArn string) string {
	out, err := regex.ParseTaskDefinitionArn(taskDefinitionArn)
	if err != nil {
		fmt.Println(err)
		return "ERROR"
	}

	return out
}

func iconizeBool(b bool) string {
	if b {
		return "[✓](fg-green)"
	}

	return "[x](fg-red)"
}

func colorizeContainerInstanceStatus(s string) string {
	switch s {
	case "ACTIVE":
		return "[ACTIVE](fg-green)"
	case "DRAINING":
		return "[DRAINING](fg-yellow)"
	}

	return s
}

func getIntEcsResource(resources []*ecs.Resource, name string) int64 {
	for _, r := range resources {
		if *r.Name == name {
			return *r.IntegerValue
		}
	}

	return 0
}

func colorizeTaskStatus(s string) string {
	switch s {
	case "RUNNING":
		return "[RUNNING](fg-green)"
	case "PENDING":
		return "[PENDING](fg-yellow)"
	case "STOPPED":
		return "[STOPPED](fg-red)"
	}

	return s
}
