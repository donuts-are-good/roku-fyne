package main

import (
	"net/http"
	"os"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	if len(os.Args) != 2 {
		println("Usage: go run main.go <roku_ip_address>")
		os.Exit(1)
	}

	rokuIP := os.Args[1]

	// Define a function that sends a keypress command to the Roku
	sendCommand := func(key string) {
		url := "http://" + rokuIP + ":8060/keypress/" + key
		resp, err := http.Post(url, "application/x-www-form-urlencoded", nil)
		if err != nil {
			println("Error:", err.Error())
			return
		}
		resp.Body.Close()
	}

	// Create a new Fyne application
	a := app.New()
	w := a.NewWindow("Roku Remote")

	// Create buttons for each Roku command
	backBtn := widget.NewButton("Back", func() { sendCommand("Back") })
	homeBtn := widget.NewButton("Home", func() { sendCommand("Home") })
	upBtn := widget.NewButton("Up", func() { sendCommand("Up") })
	downBtn := widget.NewButton("Down", func() { sendCommand("Down") })
	leftBtn := widget.NewButton("Left", func() { sendCommand("Left") })
	rightBtn := widget.NewButton("Right", func() { sendCommand("Right") })
	selectBtn := widget.NewButton("OK", func() { sendCommand("Select") })
	replayBtn := widget.NewButton("Replay", func() { sendCommand("Rev") })
	optionBtn := widget.NewButton("Option", func() { sendCommand("Info") })
	rewBtn := widget.NewButton("Rewind", func() { sendCommand("Rev") })
	playBtn := widget.NewButton("Play/Pause", func() { sendCommand("Play") })
	fwdBtn := widget.NewButton("Forward", func() { sendCommand("Fwd") })

	// Create layout and add buttons
	layout := container.NewVBox(
		container.NewGridWithColumns(2, container.NewMax(backBtn), container.NewMax(homeBtn)),
		container.NewGridWithColumns(1, container.NewMax(upBtn)),
		container.NewGridWithColumns(3, container.NewMax(leftBtn), container.NewMax(selectBtn), container.NewMax(rightBtn)),
		container.NewGridWithColumns(1, container.NewMax(downBtn)),
		container.NewGridWithColumns(2, container.NewMax(replayBtn), container.NewMax(optionBtn)),
		container.NewGridWithColumns(3, container.NewMax(rewBtn), container.NewMax(playBtn), container.NewMax(fwdBtn)),
	)
	w.SetContent(layout)

	// Show and run the application
	w.ShowAndRun()
}
