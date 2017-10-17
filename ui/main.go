package ui

import (
	"fmt"
	"strconv"

	ui "github.com/gizak/termui"

	"github.com/onoffleftright/ecs-roll/pollster"
	"github.com/onoffleftright/ecs-roll/regex"
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
	// instanceTable.FgColor = ui.ColorYellow
	instanceTable.BorderLabel = "Cluster Instances"
	instanceTable.Separator = false
	instanceTable.Height = 10
	go func() {
		c := pollster.GetContainerInstancesChannel("microservices_20170925")
		for {
			containerInstances := <-c

			rows := [][]string{
				[]string{"Container Instance", "EC2 Instance", "Agent Connected", "Status", "Running Tasks"},
			}

			for _, containerInstance := range containerInstances {
				rows = append(rows, []string{
					parseContainerInstanceId(*containerInstance.ContainerInstanceArn),
					*containerInstance.Ec2InstanceId,
					iconizeBool(*containerInstance.AgentConnected),
					colorizeStatus(*containerInstance.Status),
					strconv.Itoa(int(*containerInstance.RunningTasksCount)),
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

func parseContainerInstanceId(containterInstanceArn string) string {
	out, err := regex.ParseContainerInstanceId(containterInstanceArn)
	if err != nil {
		fmt.Println(err)
		return "ERROR"
	}

	return out
}

func iconizeBool(b bool) string {
	if b {
		return "✓"
	}

	return "x"
}

func colorizeStatus(s string) string {
	switch s {
	case "ACTIVE":
		return "[ACTIVE](fg-green)"
	case "DRAINING":
		return "[DRAINING](fg-yellow)"
	}

	return s
}
