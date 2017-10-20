package regex

import (
	"fmt"
	"regexp"
)

var containerInstanceIdFromArnRegex *regexp.Regexp
var taskArnRegex *regexp.Regexp
var taskDefinitionArnRegex *regexp.Regexp

func init() {
	containerInstanceIdFromArnRegex = regexp.MustCompile(`arn:aws:ecs:[a-z0-9-]*:\d*:container-instance\/([0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12})`)
	taskArnRegex = regexp.MustCompile(`arn:aws:ecs:[a-z0-9-]*:\d*:task\/([0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12})`)
	taskDefinitionArnRegex = regexp.MustCompile(`arn:aws:ecs:[a-z0-9-]*:\d*:task-definition/([a-z-_]*:\d*)`)
}

func ParseContainerInstanceId(containerInstanceArn string) (string, error) {
	matches := containerInstanceIdFromArnRegex.FindStringSubmatch(containerInstanceArn)

	if len(matches) != 2 {
		return "", fmt.Errorf("received malformed container instance ARN: %s", containerInstanceArn)
	}

	return matches[1], nil
}

func ParseTaskArn(taskArn string) (string, error) {
	matches := taskArnRegex.FindStringSubmatch(taskArn)

	if len(matches) != 2 {
		return "", fmt.Errorf("received malformed task arn: %s", taskArn)
	}

	return matches[1], nil
}

func ParseTaskDefinitionArn(arn string) (string, error) {
	matches := taskDefinitionArnRegex.FindStringSubmatch(arn)

	if len(matches) != 2 {
		return "", fmt.Errorf("received malformed task definition arn: %s", arn)
	}

	return matches[1], nil
}
