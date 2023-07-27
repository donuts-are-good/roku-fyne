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
	selectBtn := widget.NewButton("Select", func() { sendCommand("Select") })
	replayBtn := widget.NewButton("Replay", func() { sendCommand("Rev") })
	fwdBtn := widget.NewButton("Forward", func() { sendCommand("Fwd") })
	playBtn := widget.NewButton("Play", func() { sendCommand("Play") })

	// Create layout and add buttons
	layout := container.NewVBox(
		container.NewHBox(backBtn, homeBtn),
		container.NewHBox(upBtn),
		container.NewHBox(leftBtn, selectBtn, rightBtn),
		container.NewHBox(downBtn),
		container.NewHBox(replayBtn),
		container.NewHBox(fwdBtn),
		container.NewHBox(playBtn),
	)
	w.SetContent(layout)

	// Show and run the application
	w.ShowAndRun()
}
