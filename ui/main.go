package ui

import (
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/service/ecs"
	ui "github.com/gizak/termui"
)

const logo = `[8888888888  .d8888b.   .d8888b.](fg-yellow,fg-bold)  [                 888 888](fg-yellow)
[888        d88P  Y88b d88P  Y88b](fg-yellow,fg-bold) [                 888 888](fg-yellow)
[888        888    888 Y88b.     ](fg-yellow,fg-bold) [                 888 888](fg-yellow)
[8888888    888         "Y888b.  ](fg-yellow,fg-bold) [888d888  .d88b.  888 888](fg-yellow)
[888        888            "Y88b.](fg-yellow,fg-bold) [888P"   d88""88b 888 888](fg-yellow)
[888        888    888       "888](fg-yellow,fg-bold) [888     888  888 888 888](fg-yellow)
[888        Y88b  d88P Y88b  d88P](fg-yellow,fg-bold) [888     Y88..88P 888 888](fg-yellow)
[8888888888  "Y8888P"   "Y8888P" ](fg-yellow,fg-bold) [888      "Y88P"  888 888](fg-yellow)`

// Container Instance view members
var currentRowSelection int
var containerInstanceView *ui.Table
var containerInstances []ecs.ContainerInstance
var currentRowSelectionChangeHandler func(c ecs.ContainerInstance)
var toggleDrainContainerInstanceHandler func(c ecs.ContainerInstance)

// Task view members
var taskView *ui.Table
var tasks []ecs.Task

func init() {
	err := ui.Init()
	if err != nil {
		panic(err)
	}

	banner := ui.NewPar(logo)
	banner.Width = getWidth(logo)
	banner.Height = getHeight(logo)
	banner.Border = false

	containerInstanceView = ui.NewTable()
	containerInstanceView.BorderLabel = "Cluster Instances"
	containerInstanceView.Separator = false
	containerInstanceView.Height = 15

	taskView = ui.NewTable()
	taskView.BorderLabel = "Tasks"
	taskView.Separator = false
	taskView.Height = 15

	ui.Body.AddRows(
		ui.NewRow(
			ui.NewCol(12, 0, banner),
		),
		ui.NewRow(
			ui.NewCol(12, 0, containerInstanceView),
		),
		ui.NewRow(
			ui.NewCol(12, 0, taskView),
		),
	)

	ui.Body.Align()
	ui.Render(ui.Body)

	// Dummy handler if not registered
	currentRowSelectionChangeHandler = func(c ecs.ContainerInstance) {}
	toggleDrainContainerInstanceHandler = func(c ecs.ContainerInstance) {}

	ui.Handle("/sys/kbd/q", func(ui.Event) {
		ui.StopLoop()
	})

	ui.Handle("/sys/kbd/<up>", func(ui.Event) {
		if currentRowSelection > 0 {
			currentRowSelection--
			renderContainerInstancesView()
			currentRowSelectionChangeHandler(containerInstances[currentRowSelection])
		}
	})

	ui.Handle("/sys/kbd/<down>", func(ui.Event) {
		if currentRowSelection < len(containerInstances)-1 {
			currentRowSelection++
			renderContainerInstancesView()
			currentRowSelectionChangeHandler(containerInstances[currentRowSelection])
		}
	})

	ui.Handle("/sys/kbd/d", func(ui.Event) {
		toggleDrainContainerInstanceHandler(containerInstances[currentRowSelection])
	})
}

func Run() {
	ui.Loop()
	defer ui.Close()
}

func Stop() {
	ui.StopLoop()
}

// Container Instance View funcs

func renderContainerInstancesView() {
	rows := [][]string{
		[]string{
			"[Container Instance](fg-bold)",
			"[EC2 Instance](fg-bold)",
			"[Agent Connected](fg-bold)",
			"[Status](fg-bold)",
			"[Running Tasks](fg-bold)",
			"[CPU Available](fg-bold)",
			"[Mem Available](fg-bold)",
			"[Agent Version](fg-bold)",
			"[Docker Version](fg-bold)",
			"[Registered at](fg-bold)",
		},
	}

	for _, containerInstance := range containerInstances {
		rows = append(rows, []string{
			parseContainerInstanceId(*containerInstance.ContainerInstanceArn),
			*containerInstance.Ec2InstanceId,
			iconizeBool(*containerInstance.AgentConnected),
			colorizeContainerInstanceStatus(*containerInstance.Status),
			strconv.Itoa(int(*containerInstance.RunningTasksCount)),
			strconv.Itoa(int(getIntEcsResource(containerInstance.RemainingResources, "CPU"))),
			strconv.Itoa(int(getIntEcsResource(containerInstance.RemainingResources, "MEMORY"))),
			*containerInstance.VersionInfo.AgentVersion,
			strings.Replace(*containerInstance.VersionInfo.DockerVersion, "DockerVersion: ", "", -1),
			containerInstance.RegisteredAt.Format(time.RFC1123),
		})
	}

	containerInstanceView.Rows = rows
	containerInstanceView.Analysis()
	for i, _ := range containerInstanceView.BgColors {
		containerInstanceView.BgColors[i] = ui.ColorBlack
	}
	containerInstanceView.BgColors[currentRowSelection+1] = ui.ColorMagenta

	ui.Render(ui.Body)
}

func UpdateContainerInstances(c []ecs.ContainerInstance) {
	containerInstances = c

	if currentRowSelection > len(containerInstances)-1 {
		currentRowSelection = len(containerInstances) - 1
	}

	renderContainerInstancesView()
}

func HandleCurrentRowSelectionChange(f func(c ecs.ContainerInstance)) {
	currentRowSelectionChangeHandler = f
}

func HandleToggleDrainContainerInstance(f func(c ecs.ContainerInstance)) {
	toggleDrainContainerInstanceHandler = f
}

// Individual Container Instance Detail View

func renderContainerInstanceDetailView() {
	rows := [][]string{
		[]string{
			"[Task](fg-bold)",
			"[Task Definition](fg-bold)",
			"[Last Status](fg-bold)",
			"[Desired Status](fg-bold)",
		},
	}

	for _, t := range tasks {
		rows = append(rows, []string{
			parseTaskArn(*t.TaskArn),
			parseTaskDefinitionArn(*t.TaskDefinitionArn),
			colorizeTaskStatus(*t.LastStatus),
			*t.DesiredStatus,
		})
	}

	taskView.Rows = rows
	ui.Render(ui.Body)
}

func UpdateTasks(t []ecs.Task) {
	tasks = t
	renderContainerInstanceDetailView()
}
