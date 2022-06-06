package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/common-nighthawk/go-figure"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"gopkg.in/yaml.v2"
)

var tablePadding int = 2

var forceEnd = make(chan int)

func trak() {

	if len(flag.Args()) < 2 {
		handle(errors.New(""), "No valid trak arg was given")
		flag.Usage()
	}

	var trakArg = flag.Args()[1]

	switch trakArg {
	case "alljobs":
		liveFeed()
	case "job":
		liveFeed()
	case "list":
		listJobs()
	case "logs":
		apiLog()
	default:
		flag.Usage()
		handle(errors.New(""), "No valid trak arg was given")
	}

}

// Get log file from job
func apiLog() {
	if len(flag.Args()) < 3 {
		flag.Usage()
		handle(errors.New(""), "No valid location was given")
	}

	var location = flag.Args()[2]
	var jobID = flag.Args()[3]

	// Get sharpcd file
	f, err := ioutil.ReadFile("./sharpcd.yml")
	var con config
	err = yaml.Unmarshal(f, &con)
	handle(err, "Failed to read and extract sharpcd.yml")

	urlJobs := con.Trak[location] + "/api/logs/" + jobID

	// Insert needed data
	secret := getSec()

	apiOutput, _ := getAPIOutput(urlJobs, secret)
	logMsg := apiOutput.Message

	fmt.Println("\nLogs from " + con.Trak[location] + ":\n")
	fmt.Println(logMsg + "\n")
}

// Lists all jobs running on server
func listJobs() {
	if len(flag.Args()) < 3 {
		flag.Usage()
		handle(errors.New(""), "No valid location was given")
	}

	var location = flag.Args()[2]

	var jobIDs []string

	// Get sharpcd file
	f, err := ioutil.ReadFile("./sharpcd.yml")
	var con config
	err = yaml.Unmarshal(f, &con)
	handle(err, "Failed to read and extract sharpcd.yml")

	urlJobs := con.Trak[location] + "/api/jobs/"

	// Insert needed data
	secret := getSec()

	apiOutput, _ := getAPIOutput(urlJobs, secret)
	jobs := apiOutput.Jobs

	for _, job := range jobs {
		jobIDs = append(jobIDs, "  - "+job.ID)
	}

	fmt.Println("\nList of SharpCD Jobs running on " + con.Trak[location] + ":\n")
	fmt.Println(strings.Join(jobIDs, "\n") + "\n")
}

// Creates Interface
func liveFeed() {

	var urlJob string
	var urlJobs string
	var urlLog string
	var oldWidth = 0

	var trakAll string = "alljobs"
	var trakOne string = "job"

	if len(flag.Args()) < 3 {
		flag.Usage()
		handle(errors.New(""), "No valid location was given")
	}

	var trakArg = flag.Args()[1]
	var location = flag.Args()[2]

	// Get sharpcd file
	f, err := ioutil.ReadFile("./sharpcd.yml")
	var con config
	err = yaml.Unmarshal(f, &con)
	handle(err, "Failed to read and extract sharpcd.yml")

	urlJobs = con.Trak[location] + "/api/jobs/"

	// Insert needed data
	secret := getSec()

	// Only needed for single job requests
	if trakArg == trakOne {
		if len(flag.Args()) < 4 {
			flag.Usage()
			handle(errors.New(""), "No valid Job ID was given")
		}

		var jobID = flag.Args()[3]

		urlJob = con.Trak[location] + "/api/job/" + jobID
		urlLog = con.Trak[location] + "/api/logsfeed/" + jobID

		// Tests to check if job exists
		apiOutput, code := getAPIOutput(urlJob, secret)

		if code == statCode.Accepted {
			fmt.Printf("Connection to API...")
		} else {
			fmt.Println(apiOutput.Message)
			fmt.Println()
			os.Exit(1)
		}
	}

	// Tests to ensure you can actually reach the server
	apiOutput, code := getAPIOutput(urlJobs, secret)

	if code == statCode.Accepted {
		fmt.Printf("Connection to API...")
	} else {
		fmt.Println(apiOutput.Message)
		fmt.Println()
		fmt.Printf("API Connection Failed!\n")
		os.Exit(1)
	}

	// Load UI
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	// Function for rendering UI
	draw := func() {
		table := tableMain()

		var jobs []*taskJob

		// Gets all rows if there is one job
		if trakArg == trakOne {
			apiOutput, _ := getAPIOutput(urlJob, secret)
			job := apiOutput.Job

			jobs = append(jobs, job)
		}

		// Gets all rows if there is there are multiple jobs
		if trakArg == trakAll {
			apiOutput, _ = getAPIOutput(urlJobs, secret)
			jobs = apiOutput.Jobs
		}

		// Adds in a row for each job
		for _, job := range jobs {
			jobStr := []string{
				job.ID, job.Name, job.Type, job.Status, job.ErrMsg, strconv.FormatBool(job.Reconnect),
			}
			table.Rows = append(table.Rows, jobStr)
		}

		// Makes sure all widths are correct
		table.ColumnWidths = generateColumnWidths(table.Rows, table.Rows[0])
		width := sum(table.ColumnWidths) + len(table.Rows[0]) + 1

		// Deals with bug that if width space shrinks it doesn't render correctly
		if oldWidth > width {
			width = oldWidth
		} else {
			oldWidth = width
		}

		title, heightTitle := createTitle(width)
		ui.Render(title)

		// Creates Table
		deltaYTable := heightTitle
		heightTable := deltaYTable + len(table.Rows)*2 + 1
		table.SetRect(0, deltaYTable, width, heightTable)
		ui.Render(table)

		// Creates Logs if needed
		var heightLogs int
		var trakLogs *widgets.Paragraph

		if trakArg == trakOne {
			trakLogs, heightLogs = logging(width, heightTable, urlLog, secret)
			ui.Render(trakLogs)
		} else {
			heightLogs = heightTable
		}

		close := closing(width, tablePadding, heightLogs)
		ui.Render(close)

	}

	// Runs Interface
	uiEvents := ui.PollEvents()
	ticker := time.NewTicker(time.Second).C
	draw()

	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				return
			}
		case <-ticker:
			draw()

		}
	}

}

// Creates Table with styling
func tableMain() *widgets.Table {

	headers := []string{"ID", "Name", "Type", "Status", "Error Message", "Reconnected"}

	table1 := widgets.NewTable()

	table1.TextStyle = ui.NewStyle(ui.ColorWhite)
	table1.TextAlignment = ui.AlignCenter
	table1.FillRow = true

	table1.Rows = [][]string{headers}
	table1.BorderStyle.Fg = ui.ColorCyan
	table1.RowStyles[0] = ui.NewStyle(ui.ColorGreen, ui.ColorClear, ui.ModifierBold)

	return table1
}

// Creates the fancy title font
func createTitle(width int) (*widgets.Paragraph, int) {
	title := widgets.NewParagraph()
	myFigure := figure.NewFigure("SharpCD Trak", "", true)
	figureRows := myFigure.Slicify()
	title.Text = ""
	titleSpace := 0

	// Centers the Figure
	if width < len(figureRows[0]) {
		width = len(figureRows[0]) + 4
		titleSpace = 0
	} else {
		titleSpace = (width-len(figureRows[0]))/2 - 2*tablePadding - 1
	}

	// Check to make sure there is room
	if titleSpace < 1 {
		titleSpace = 0
	}

	// Draw the Figure into a string
	for _, row := range figureRows {
		title.Text += (strings.Repeat(" ", titleSpace) + row + "\n")
	}

	heightTitle := len(figureRows) + tablePadding
	title.SetRect(0, 0, width, heightTitle)
	title.BorderStyle.Fg = ui.ColorCyan
	title.TextStyle = ui.NewStyle(ui.ColorGreen, ui.ColorClear, ui.ModifierBold)

	return title, heightTitle
}

// Creates closing screen
func closing(width int, tablePadding int, heightLogs int) *widgets.Paragraph {

	deltaYClose := heightLogs
	heightClose := deltaYClose + 3

	close := widgets.NewParagraph()

	dt := time.Now()

	text1 := " Press Ctrl+C to Exit"
	text2 := dt.Format("01-02-2006 15:04:05") + " "
	space := width - len(text1) - len(text2) - tablePadding

	// Ensure the text is aligned and padded
	close.Text = text1 + strings.Repeat(" ", space) + text2
	close.SetRect(0, deltaYClose, width, heightClose)
	close.BorderStyle.Fg = ui.ColorCyan

	close.TextStyle = ui.NewStyle(ui.ColorGreen, ui.ColorClear, ui.ModifierBold)

	return close
}

// Creates the logging page
func logging(width int, heightTable int, urlLog string, secret string) (*widgets.Paragraph, int) {
	logs, _ := getAPIOutput(urlLog, secret)
	msg := logs.Message

	deltaYLogs := heightTable
	heightLogs := deltaYLogs + 22

	trakLogs := widgets.NewParagraph()
	trakLogs.Text = "\n" + msg
	trakLogs.SetRect(0, deltaYLogs, width, heightLogs)

	trakLogs.BorderStyle.Fg = ui.ColorCyan
	trakLogs.Title = "Live Output"

	trakLogs.TitleStyle = ui.NewStyle(ui.ColorGreen, ui.ColorClear, ui.ModifierBold)

	return trakLogs, heightLogs
}

// Finds the sum of ints
func sum(array []int) int {
	result := 0
	for _, v := range array {
		result += v
	}
	return result
}

// Finds the longest string in column (index)
func longest(array [][]string, index int) int {
	longestString := 0

	for _, row := range array {
		if new := len(row[index]); new > longestString {
			longestString = new
		}
	}

	if longestString > 40 {
		longestString = 40
	}

	return longestString + 4
}

func generateColumnWidths(rows [][]string, columns []string) []int {

	columnWidths := []int{}

	for index := range columns {
		columnWidths = append(columnWidths, longest(rows, index))
	}

	return columnWidths
}

func getAPIOutput(url string, secret string) (response, int) {

	trakPayload := trakPostData{
		Secret:  secret,
		Version: sharpCDVersion,
	}

	// Create request
	jsonStr, err := json.Marshal(trakPayload)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	handle(err, "Failed to create request for"+url)

	// Create client
	// Allow self-signed certs
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	// Do Request
	resp, err := client.Do(req)
	handle(err, "Failed to connect to SharpCD Trak: "+url)
	defer resp.Body.Close()

	// Read Body and Status
	body, err := ioutil.ReadAll(resp.Body)
	handle(err, "Failed to read body of response"+url)

	var apiOutput response
	err = json.Unmarshal(body, &apiOutput)
	handle(err, "Failed to unmarshal body"+url)

	return apiOutput, resp.StatusCode
}
