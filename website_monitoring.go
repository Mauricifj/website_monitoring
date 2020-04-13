package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// Will check all urls this number of times
const rounds = 5

// Will wait for this amount of time between rounds
const delay = 5
const timeUnit = time.Second

func main() {
	welcome()

	for {
		showMenu()

		option := readOption()

		switch option {
		case 1:
			startMonitoring()

		case 2:
			extractLogsFromTxt()
			newLine()

		case 0:
			fmt.Println("Process finishing...")
			os.Exit(0)

		default:
			fmt.Println("Invalid option")
			os.Exit(-1)
		}
	}
}

// Display and menu
func welcome() {
	version := "1.3"
	fmt.Println("*** Program Website Monitoring - Version", version, "***")
}

func showMenu() {
	fmt.Println("--- MENU ---")
	fmt.Println("1 - Start monitoring")
	fmt.Println("2 - Show logs")
	fmt.Println("0 - Quit")
}

func readOption() int {
	fmt.Print("Choose one option, please: ")
	var option int
	fmt.Scan(&option)
	return option
}

func newLine() {
	fmt.Println()
}

// Monitoring
func startMonitoring() {
	fmt.Println("-- WEBSITE MONITOR STARTED --")
	newLine()

	urls := extractUrlsFromTxt()

	for i := 0; i < rounds; i++ {
		for _, url := range urls {
			check(url)
		}

		// Check if it's the last loop, if it's not, sleep for specified time
		if i != rounds-1 {
			newLine()
			fmt.Println("-- WAIT FOR THE NEXT MONITORING --")
			newLine()
			time.Sleep(delay * timeUnit)
		} else {
			newLine()
			fmt.Println("-- WEBSITE MONITOR ENDED --")
			newLine()
		}
	}
}

func check(url string) {
	response, err := http.Get(url)

	if err != nil {
		fmt.Println("Website:", url, " - Error:", err)
	} else {
		if response.StatusCode == 200 {
			fmt.Println("Url:", url, "is ok")
			log(url, true)
		} else {
			fmt.Println("Url:", url, "is not ok. Status code:", response.StatusCode)
			log(url, false)
		}
	}
}

// Files
func extractUrlsFromTxt() []string {
	var urls []string

	fileName := "urls.txt"

	file, err := os.Open(fileName)

	if err != nil {
		fmt.Println("Error on ", fileName, ":", err)
	} else {
		reader := bufio.NewReader(file)

		for {
			url, err := reader.ReadString('\n')
			url = strings.TrimSpace(url)
			urls = append(urls, url)

			if err == io.EOF {
				break
			}
		}

		file.Close()
	}
	return urls
}

func log(url string, status bool) {
	fileName := "log.txt"
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		fmt.Println("Error on opening/creating", fileName, ":", err)
	} else {
		file.WriteString(time.Now().Format("02/01/2006 15:04:05") + " - url:" + url + " - online:" + strconv.FormatBool(status) + "\n")
	}

	file.Close()
}

func extractLogsFromTxt() {
	fileName := "log.txt"
	file, err := ioutil.ReadFile(fileName)

	if err != nil {
		fmt.Println("Error on opening", fileName, ":", err)
	} else {
		fmt.Println(string(file))
	}
}
