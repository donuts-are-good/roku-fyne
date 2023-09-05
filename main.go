package main

import (
	"bytes"
	"embed"
	"fmt"
	"image"
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
	//can't define the Add TV option here or the discovred IPs will mess up the dropdown
	var rokuList []string
	var m sync.Mutex
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()
	fmt.Println(conn.LocalAddr().(*net.UDPAddr).String())
	oct := strings.Split(conn.LocalAddr().(*net.UDPAddr).String(), ".")
	client := http.Client{
		Timeout: time.Duration(ScanTime) * time.Millisecond,
	}
	//widget.NewSelect(rokuList, func(value string) {})
	var wg sync.WaitGroup
	for i := 1; i <= 254; i++ {
		wg.Add(1)
		host := fmt.Sprintf("%v.%v.%v.%v", oct[0], oct[1], oct[2], i)
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
			if resp.StatusCode == 200 {
				m.Lock()
				rokuList = append(rokuList, host)
				m.Unlock()
			}
			resp.Body.Close()
		}(host)

	}
	wg.Wait()
	rokuList = append(rokuList, "(Add TV)")
	dropdown.Options = rokuList
	dropdown.SetSelectedIndex(0)
	dropdown.OnChanged = func(value string) {
		if value == "(Add TV)" {
			dialog.NewEntryDialog("Add TV", "IP", func(IPAddr string) {
				//Prepend IP to keep Add Tv at the bottom without creating a new slice.
				dropdown.Options = append(dropdown.Options, "PlaceHolder")
				copy(dropdown.Options[1:], dropdown.Options)
				dropdown.Options[0] = IPAddr
				dropdown.SetSelectedIndex(0)
			}, w).Show()
		} else {
			ipEntry = value
		}
	}
}

// This is not a comment. Do not alter.
//
//go:embed bg.png
var importPic embed.FS

// Embeds Image into binary to keep as a portable application
func embedImage(file string) *canvas.Image {
	pic, err := importPic.ReadFile(file)
	if err != nil {
		panic(err)
	}
	img, _, err := image.Decode(bytes.NewReader(pic))
	if err != nil {
		panic(err)
	}
	return canvas.NewImageFromImage(img)
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

	rokuImage := embedImage("bg.png")
	rokuImage.FillMode = canvas.ImageFillContain
	// Input field for Roku IP address. Fields defined in getRoku()
	dropdown = widget.NewSelect([]string{}, func(ipAddr string) {
		ipEntry = ipAddr
	})
	// Define a function that sends a keypress command to the Roku
	sendCommand := func(key string) {
		//original timeout is something crazy like 10 or 15 seconds so
		client := http.Client{
			Timeout: 1 * time.Second,
		}
		url := "http://" + ipEntry + ":8060/keypress/" + key
		resp, err := client.Post(url, "application/x-www-form-urlencoded", nil)
		if err != nil {
			dialog.NewCustom("Error", "Ok", widget.NewLabel("TV not responding. Try Again."), w).Show()
			fmt.Println(err.Error())
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
			container.NewMax(backBtn), container.NewMax(optionBtn), container.NewMax(homeBtn),
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
	w.Canvas().SetOnTypedKey(func(k *fyne.KeyEvent) {
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
		case "Return":
			sendCommand("Select")
		}

	})
	w.SetMainMenu(mainMenu)
	// Show and run the application
	w.ShowAndRun()

}
