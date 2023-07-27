package main

import (
	"net/http"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func main() {
	// Create a new Fyne application
	a := app.New()
	w := a.NewWindow("れもとく")

	// Set the window size
	w.Resize(fyne.NewSize(200, 351))
	w.SetFixedSize(true)

	bg := canvas.NewImageFromFile("bg.png")
	bg.FillMode = canvas.ImageFillStretch

	// Input field for Roku IP address
	ipEntry := widget.NewEntry()
	ipEntry.SetPlaceHolder("Enter IP address here...")

	// Define a function that sends a keypress command to the Roku
	sendCommand := func(key string) {
		url := "http://" + ipEntry.Text + ":8060/keypress/" + key
		resp, err := http.Post(url, "application/x-www-form-urlencoded", nil)
		if err != nil {
			println("Error:", err.Error())
			return
		}
		resp.Body.Close()
	}

	// Create buttons for each Roku command
	backBtn := widget.NewButton("Back", func() { sendCommand("Back") })
	homeBtn := widget.NewButton("Home", func() { sendCommand("Home") })
	upBtn := widget.NewButton("Up", func() { sendCommand("Up") })
	downBtn := widget.NewButton("Down", func() { sendCommand("Down") })
	leftBtn := widget.NewButton("Left", func() { sendCommand("Left") })
	rightBtn := widget.NewButton("Right", func() { sendCommand("Right") })
	selectBtn := widget.NewButton("OK", func() { sendCommand("Select") })
	replayBtn := widget.NewButton("Replay", func() { sendCommand("InstantReplay") })
	optionBtn := widget.NewButton("Option", func() { sendCommand("Info") })
	rewBtn := widget.NewButton("Rew", func() { sendCommand("Rev") })
	playBtn := widget.NewButton("Play", func() { sendCommand("Play") })
	fwdBtn := widget.NewButton("Fwd", func() { sendCommand("Fwd") })

	// Create layout and add buttons
	controls := container.NewVBox(
		ipEntry,
		container.NewGridWithColumns(2, container.NewMax(backBtn), container.NewMax(homeBtn)),
		container.NewGridWithColumns(1, container.NewMax(upBtn)),
		container.NewGridWithColumns(3, container.NewMax(leftBtn), container.NewMax(selectBtn), container.NewMax(rightBtn)),
		container.NewGridWithColumns(1, container.NewMax(downBtn)),
		container.NewGridWithColumns(2, container.NewMax(replayBtn), container.NewMax(optionBtn)),
		container.NewGridWithColumns(3, container.NewMax(rewBtn), container.NewMax(playBtn), container.NewMax(fwdBtn)),
	)

	// Add the background and the controls to the window content
	w.SetContent(container.New(layout.NewMaxLayout(), bg, controls))

	// Show and run the application
	w.ShowAndRun()
}
