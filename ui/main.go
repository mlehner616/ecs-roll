package ui

import (
	ui "github.com/gizak/termui"
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

	ui.Render(banner)

	ui.Handle("/sys/kbd/q", func(ui.Event) {
		ui.StopLoop()
	})

	ui.Loop()
}
