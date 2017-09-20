package ui

import (
	ui "github.com/gizak/termui"
	aws "github.com/onoffleftright/ecs-roll/pollster"
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

	instanceList := ui.NewList()
	instanceList.ItemFgColor = ui.ColorYellow
	instanceList.BorderLabel = "Cluster Instances"
	instanceList.Height = 10
	go func() {
		c := aws.GetContainerInstancesChannel("blackbird_microservices")
		for {
			instanceList.Items = <-c
			ui.Render(ui.Body)
		}
	}()

	ui.Body.AddRows(
		ui.NewRow(
			ui.NewCol(6, 0, banner),
		),
		ui.NewRow(
			ui.NewCol(6, 0, instanceList),
		),
	)

	ui.Body.Align()
	ui.Render(ui.Body)

	ui.Handle("/sys/kbd/q", func(ui.Event) {
		ui.StopLoop()
	})

	ui.Loop()
}
