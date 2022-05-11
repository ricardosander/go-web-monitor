package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

const monitoringTimes = 3
const monitoringDelayInSeconds = 5

func main() {

	showIntroduction()

	for {

		showMenu()
		command := readCommand()

		switch command {
		case 1:
			startMonitoring()
		case 2:
			showLogs()
		case 0:
			exitProgram()
		default:
			handleInvalidCommand()
		}
	}
}

func showIntroduction() {
	var name string = "Lopes"
	version := 1.1
	fmt.Println("Hello mr.", name)
	fmt.Println("This software is on version", version)
}

func showMenu() {
	fmt.Println("1- Lauch monitoring")
	fmt.Println("2- Show logs")
	fmt.Println("0- Exit")
}

func readCommand() int {
	var command int
	fmt.Scan(&command)
	fmt.Println("Choose command:", command)
	return command
}

func exitProgram() {
	fmt.Println("Exiting program")
	os.Exit(0)
}

func handleInvalidCommand() {
	fmt.Println("Invalid command")
	os.Exit(-1)
}

func startMonitoring() {
	fmt.Println("Monitoring...")

	var sites = readSitesFromFile()
	var total = len(sites)
	s, _ := json.Marshal(sites)
	fmt.Println(string(s))

	for i := 0; i < monitoringTimes; i++ {
		for index, site := range sites {
			fmt.Println("Checking sites:", index+1, "/", total)
			checkSite(site)
		}
		fmt.Println("Sleeping to check again")
		time.Sleep(monitoringDelayInSeconds * time.Second)
	}
	fmt.Println("Finishing monitoring")
}

func checkSite(site string) {
	resp, error := http.Get(site)

	if error != nil {
		fmt.Println("Error when trying to reach ", site, ":", error)
		return
	}

	if resp != nil && resp.StatusCode == 200 {
		fmt.Println(site, "is up :)")
		writeLog(site, "up")
		return
	}

	fmt.Println(site, "is down :( Status Code:", resp.StatusCode)
	writeLog(site, "down")
}

func readSitesFromFile() []string {
	var file, error = os.Open("sites.txt")
	if error != nil {
		fmt.Println("File cannot be read:", error)
		return []string{}
	}

	var reader = bufio.NewReader(file)
	var sites = []string{}
	for {

		line, error := reader.ReadString('\n')

		if error != nil && error != io.EOF {
			fmt.Println("Error when readind file", file.Name(), ":", error)
			break
		}

		sites = append(sites, strings.TrimSpace(line))

		if error == io.EOF {
			break
		}
	}
	file.Close()
	return sites
}

func writeLog(site string, status string) {
	file, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		fmt.Println("Error", err)
		return
	}

	now := time.Now().Format("02/01/2006 15:04:05")

	file.WriteString(now + " - " + site + " is " + status + "\n")
	file.Close()
}

func showLogs() {
	fmt.Println("Showing logs")

	var file, error = ioutil.ReadFile("log.txt")
	if error != nil {
		fmt.Println("File cannot be read:", error)
		return
	}

	fmt.Println(string(file))
}
