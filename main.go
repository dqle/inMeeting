package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"runtime"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

func main() {
	app := app.New()
	window := app.NewWindow("inMeeting")

	//-----------------------------------------
	// Setup Systray
	//-----------------------------------------
	if desk, ok := app.(desktop.App); ok {
		systemTray := fyne.NewMenu("inMeeting",
			fyne.NewMenuItem("Start", func() { log.Println("Tapped show") }),
			fyne.NewMenuItem("Settings", func() {
				window.Show()
			}))
		desk.SetSystemTrayMenu(systemTray)
	}

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
	window.ShowAndRun()
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

	// Return "" if file doens't exist
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

	// Return the ip address
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
