package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"

	mediaDevices "github.com/dqle/go-media-devices-state"
)

func main() {
	app := app.New()
	window := app.NewWindow("inMeeting")

	window.SetIcon(appIcon)

	//-----------------------------------------
	// Setup Systray
	//-----------------------------------------
	desk, _ := app.(desktop.App)
	systemTray := fyne.NewMenu("inMeeting",
		fyne.NewMenuItem("Settings", func() {
			window.Show()
		}))
	desk.SetSystemTrayMenu(systemTray)
	desk.SetSystemTrayIcon(offIcon)

	//-----------------------------------------
	// Setup Settings Windows
	//-----------------------------------------
	ipAddrEntry := widget.NewEntry()
	ipAddrEntry.SetPlaceHolder("Enter Pi Zero IP Address")

	// Load IP address from settings.cfg if it exists
	ipAddr := loadFromFile()
	if ipAddr != "" {
		ipAddrEntry.SetText(ipAddr)
	}

	saveButton := widget.NewButton("Save", func() {
		ipAddr := ipAddrEntry.Text
		saveToFile(ipAddr)
		window.Hide()
	})

	window.SetContent(container.NewVBox(
		widget.NewLabel("Enter IP Address:"),
		ipAddrEntry,
		saveButton,
	))

	window.Resize(fyne.NewSize(300, 100))
	window.SetCloseIntercept(func() {
		window.Hide()
	})

	//-----------------------------------------
	// Loop to watch device states in a separate goroutine
	//-----------------------------------------
	go func() {
		for {
			ipAddr := loadFromFile()
			isCameraOn, cameraErr := mediaDevices.IsCameraOn()
			isMicrophoneOn, microphoneErr := mediaDevices.IsMicrophoneOn()
			if cameraErr != nil {
				log.Println(cameraErr)
			} else if microphoneErr != nil {
				log.Println(microphoneErr)
			} else {
				if isCameraOn || isMicrophoneOn {
					response, err := http.Get("http://" + ipAddr + "/status")
					if err != nil {
						log.Println(err)
						continue
					}

					statusBody, err := io.ReadAll(response.Body)
					if err != nil {
						log.Println(err)
						continue
					}
					if string(statusBody) == "off" {
						_, err := http.Post("http://"+ipAddr+"/api/on", "", nil)
						if err != nil {
							log.Println(err)
						}
						desk.SetSystemTrayIcon(appIcon)
						systemTray.Refresh()
					}
					log.Println("in meeting")

				} else {
					response, err := http.Get("http://" + ipAddr + "/status")
					if err != nil {
						log.Println(err)
						continue
					}

					statusBody, err := io.ReadAll(response.Body)
					if err != nil {
						log.Println(err)
						continue
					}

					if string(statusBody) == "on" {
						_, err := http.Post("http://"+ipAddr+"/api/off", "", nil)
						if err != nil {
							log.Println(err)
						}
						desk.SetSystemTrayIcon(offIcon)
						systemTray.Refresh()
					}
					log.Println("not in meeting")
				}
			}
			time.Sleep(3 * time.Second)
		}
	}()

	app.Run()
}

// -----------------------------------------
// Function:    settingsFileDir
// Description: Retrieve the directory of the settings.cfg file
// -----------------------------------------

func settingsFileDir() (string, string) {
	dirPath := ""
	filePath := ""

	if runtime.GOOS == "windows" {
		dirPath = filepath.Join(os.Getenv("APPDATA"), "inMeeting")
		filePath = filepath.Join(dirPath, "settings.cfg")
	} else if runtime.GOOS == "darwin" {
		dirPath = filepath.Join(os.Getenv("HOME"), "inMeeting")
		filePath = filepath.Join(dirPath, "settings.cfg")
	}

	return dirPath, filePath
}

// -----------------------------------------
// Function:    loadFromFile
// Description: Load variables in settings.cfg to the settings window
// -----------------------------------------
func loadFromFile() string {
	_, filePath := settingsFileDir()

	// Return "" if file doesn't exist
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return ""
	}

	// Return "" if there's error opening the file
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening file %s: %v\n", filePath, err)
		return ""
	}
	defer file.Close()

	// Return the IP address
	scanner := bufio.NewScanner(file)
	ipRegex := regexp.MustCompile(`ip = "(.*)"`)

	for scanner.Scan() {
		line := scanner.Text()
		if ipAddress := ipRegex.FindStringSubmatch(line); ipAddress != nil {
			return ipAddress[1]
		}
	}

	// Return null by default
	return ""
}

// -----------------------------------------
// Function:    saveToFile
// Description: Save Pi Zero IP Address to settings file
// -----------------------------------------
func saveToFile(ipAddr string) {

	dirPath, filePath := settingsFileDir()

	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		fmt.Printf("Error creating directory %s: %v\n", dirPath, err)
		return
	}

	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error creating or opening file %s: %v\n", filePath, err)
		return
	}
	defer f.Close()

	configText := fmt.Sprintf("ip = \"%s\"", ipAddr)
	_, err = f.WriteString(configText)
	if err != nil {
		fmt.Printf("Error writing to file %s: %v\n", filePath, err)
	}
}
