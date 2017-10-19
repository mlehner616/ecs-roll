package ui

import (
	"strconv"
	"strings"
	"time"

	ui "github.com/gizak/termui"

	"github.com/onoffleftright/ecs-roll/pollster"
)

const logo = `8888888888  .d8888b.   .d8888b.                   888 888
888        d88P  Y88b d88P  Y88b                  888 888
888        888    888 Y88b.                       888 888
8888888    888         "Y888b.   888d888  .d88b.  888 888
888        888            "Y88b. 888P"   d88""88b 888 888
888        888    888       "888 888     888  888 888 888
888        Y88b  d88P Y88b  d88P 888     Y88..88P 888 888
8888888888  "Y8888P"   "Y8888P"  888      "Y88P"  888 888`

func DoIt() {
	err := ui.Init()
	if err != nil {
		panic(err)
	}
	defer ui.Close()

	banner := ui.NewPar(logo)
	banner.Width = getWidth(logo)
	banner.Height = getHeight(logo)
	banner.Border = false

	instanceTable := ui.NewTable()
	instanceTable.BorderLabel = "Cluster Instances"
	instanceTable.Separator = false
	instanceTable.Height = 10
	go func() {
		c := pollster.GetContainerInstancesChannel("microservices_20170925")
		for {
			containerInstances := <-c

			rows := [][]string{
				[]string{"Container Instance", "EC2 Instance", "Agent Connected", "Status", "Running Tasks", "CPU Available", "Mem Available", "Agent Version", "Docker Version", "Registered at"},
			}

			for _, containerInstance := range containerInstances {
				rows = append(rows, []string{
					parseContainerInstanceId(*containerInstance.ContainerInstanceArn),
					*containerInstance.Ec2InstanceId,
					iconizeBool(*containerInstance.AgentConnected),
					colorizeStatus(*containerInstance.Status),
					strconv.Itoa(int(*containerInstance.RunningTasksCount)),
					strconv.Itoa(int(getIntEcsResource(containerInstance.RemainingResources, "CPU"))),
					strconv.Itoa(int(getIntEcsResource(containerInstance.RemainingResources, "MEMORY"))),
					*containerInstance.VersionInfo.AgentVersion,
					strings.Replace(*containerInstance.VersionInfo.DockerVersion, "DockerVersion: ", "", -1),
					containerInstance.RegisteredAt.Format(time.RFC1123),
				})
			}

			instanceTable.Rows = rows
			ui.Render(ui.Body)
		}
	}()

	ui.Body.AddRows(
		ui.NewRow(
			ui.NewCol(12, 0, banner),
		),
		ui.NewRow(
			ui.NewCol(12, 0, instanceTable),
		),
	)

	ui.Body.Align()
	ui.Render(ui.Body)

	ui.Handle("/sys/kbd/q", func(ui.Event) {
		ui.StopLoop()
	})

	ui.Loop()
}
