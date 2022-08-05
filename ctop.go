package main

import (
	"fmt"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"log"
	"time"
)

var loadMonitor LoadAverage

func genSummary() []string {
	currentTime := time.Now().Local().Format("15:04:05")
	tc := GetTaskCount()

	return []string{
		fmt.Sprintf("topic - %s  load average: %s", currentTime, loadMonitor.GetLoad()),
		fmt.Sprintf("Tasks: [%3d](mod:bold) total, [%3d](mod:bold) running, [%3d](mod:bold) sleeping, [%3d](mod:bold) stopped, [%3d](mod:bold) zombie",
			tc.Total, tc.Running, tc.Sleeping, tc.Stopped, tc.Zombie),
	}
}
func main() {

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()
	termWidth, termHeight := ui.TerminalDimensions()
	summary := widgets.NewList()
	summary.Rows = genSummary()
	summary.TextStyle = ui.NewStyle(ui.ColorWhite)
	summary.WrapText = false
	summary.SetRect(-1, -1, termWidth, termHeight)
	summary.Border = false

	go func() {
		loadMonitor.Run()
	}()

	uiEvents := ui.PollEvents()
	ticker := time.NewTicker(time.Second).C
	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				return
			}
		case <-ticker:
			summary.Rows = genSummary()
			//processes.Text = genProcesses(taskMonitor)
			ui.Render(summary)
		}
	}
}
