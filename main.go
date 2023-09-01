package main

import (
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func getRoku(ScanTime int) {
	var rokuList []string
	var m sync.Mutex
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()
	oct := strings.Split(conn.LocalAddr().String(), ".")
	client := http.Client{
		Timeout: time.Duration(ScanTime) * time.Millisecond,
	}
	//widget.NewSelect(rokuList, func(value string) {})
	var wg sync.WaitGroup
	for i:=1;i<=254;i++{
		wg.Add(1)
		host := fmt.Sprintf("%v.%v.%v.%v",  oct[0], oct[1], oct[2], i)
		go func(ip string) {
			defer wg.Done()
			url := "http://" + ip + ":8060/query/device-info"
			resp, err := client.Get(url)
			if err != nil {
				if strings.Contains(err.Error(), "Client.Timeout exceeded") {
					return
				}
				dialog.NewError(err, w)
			}
			if (resp.StatusCode == 200) {
				m.Lock()
				rokuList = append(rokuList, host)
				m.Unlock()
			}
			resp.Body.Close()
		}(host)
		
	}
	wg.Wait()
	dropdown.Options = rokuList
	dropdown.SetSelected(rokuList[0])
}

var w fyne.Window
var ipEntry string
var dropdown *widget.Select

func main() {
	//start looking for Roku TVs on the local network.
	go getRoku(200)
	// Create a new Fyne application
	a := app.New()
	w = a.NewWindow("Roku")
	w.Resize(fyne.NewSize(225, 350))
	rokuImage := canvas.NewImageFromFile("bg.png")
	rokuImage.FillMode = canvas.ImageFillContain
	//not used, but empty slice for filler
	var list []string
	// Input field for Roku IP address
	dropdown = widget.NewSelect(list, func(ipAddr string) {
		ipEntry = ipAddr
	})
	// Define a function that sends a keypress command to the Roku
	sendCommand := func(key string) {
		url := "http://" + ipEntry + ":8060/keypress/" + key
		resp, err := http.Post(url, "application/x-www-form-urlencoded", nil)
		if err != nil {
			println("Error:", err.Error())
			return
		}
		resp.Body.Close()
	}
	// Create buttons for each Roku command

	powerOnBtn := widget.NewButton("On", func() { sendCommand("PowerOn") })
	powerOffBtn := widget.NewButton("Off", func() { sendCommand("PowerOff") })
	backBtn := widget.NewButtonWithIcon("", theme.MailReplyIcon(), func() { sendCommand("Back") })
	homeBtn := widget.NewButtonWithIcon("", theme.HomeIcon(), func() { sendCommand("Home") })
	upBtn := widget.NewButtonWithIcon("", theme.MoveUpIcon(), func() { sendCommand("Up") })
	downBtn := widget.NewButtonWithIcon("", theme.MoveDownIcon(), func() { sendCommand("Down") })
	leftBtn := widget.NewButtonWithIcon("", theme.NavigateBackIcon(), func() { sendCommand("Left") })
	rightBtn := widget.NewButtonWithIcon("", theme.NavigateNextIcon(), func() { sendCommand("Right") })
	selectBtn := widget.NewButton("OK", func() { sendCommand("Select") })
	optionBtn := widget.NewButtonWithIcon("", theme.SettingsIcon(), func() { sendCommand("Info") })
	vDownBtn := widget.NewButtonWithIcon("", theme.VolumeDownIcon(), func() { sendCommand("VolumeDown") })
	vUpBtn := widget.NewButtonWithIcon("", theme.VolumeUpIcon(), func() { sendCommand("VolumeUP") })

	// Create layout and add buttons
	controls := container.NewVBox(
		dropdown,
		container.NewGridWithColumns(2, powerOffBtn, powerOnBtn),
		container.NewGridWithColumns(3, 
			container.NewMax(backBtn),container.NewMax(optionBtn),container.NewMax(homeBtn),
			widget.NewLabel(""), container.NewMax(upBtn), widget.NewLabel(""),
			container.NewMax(leftBtn), container.NewMax(selectBtn), container.NewMax(rightBtn),
			container.NewMax(vDownBtn), container.NewMax(downBtn), container.NewMax(vUpBtn),
		),
		//Empty label required to get image to play nice
		container.NewGridWithRows(2, widget.NewLabel(""), rokuImage),
	)

	// Add the background and the controls to the window content
	w.SetContent(container.New(layout.NewMaxLayout(), controls))

	// Create a new window for the "About" section
	aboutWindow := a.NewWindow("About")
	aboutWindow.SetContent(widget.NewLabel("This is a simple Roku remote application."))
	aboutWindow.Resize(fyne.NewSize(300, 200))

	// Create a menu
	mainMenu := fyne.NewMainMenu(
		// A quit item will be appended to our first menu
		fyne.NewMenu("File", 
			fyne.NewMenuItem("Add TV", func() { 
				dialog.NewEntryDialog("Add TV", "IP", func(IPAddr string) {
					dropdown.Options = append(dropdown.Options, IPAddr)
			}, w).Show()
			}),
			fyne.NewMenuItem("Scan for Roku", func() { 
				dropdown.Options = []string{"Please Wait..."}
				dropdown.SetSelectedIndex(0)
				dropdown.Refresh()
				getRoku(500)
		 	}),
			fyne.NewMenuItem("Quit", func() { a.Quit() }),
		),
		fyne.NewMenu("Help", fyne.NewMenuItem("About", func() {
			aboutWindow.Show()
		})),

	)
	w.Canvas().SetOnTypedKey(func (k *fyne.KeyEvent) {
		switch k.Name {
		case fyne.KeyLeft:
				sendCommand("Left")
		case fyne.KeyRight:
				sendCommand("Right")
		case fyne.KeySpace:
			sendCommand("Select")
		case fyne.KeyUp:
				sendCommand("Up")
		case fyne.KeyDown:
				sendCommand("Down")
		case fyne.KeyBackspace:
			sendCommand("Back")
		}

	})
	w.SetMainMenu(mainMenu)
	// Show and run the application
	w.ShowAndRun()
	
}
