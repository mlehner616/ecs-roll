package pollster

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
)

var svc *ecs.ECS

func init() {
	svc = ecs.New(session.New())
}
