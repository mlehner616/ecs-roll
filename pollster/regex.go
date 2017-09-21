package pollster

import (
	"fmt"
	"regexp"
)

var containerInstanceIdFromArnRegEx *regexp.Regexp

func init() {
	containerInstanceIdFromArnRegEx = regexp.MustCompile(`arn:aws:ecs:[a-z0-9-]*:\d*:container-instance\/([0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12})`)
}

func parseContainerInstanceId(containerInstanceArn string) (string, error) {
	matches := containerInstanceIdFromArnRegEx.FindStringSubmatch(containerInstanceArn)

	if len(matches) != 2 {
		return "", fmt.Errorf("received malformed container instance ARN: %s", containerInstanceArn)
	}

	return matches[1], nil
}
