package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/common-nighthawk/go-figure"
	ui "github.com/gizak/termui"
	"github.com/gizak/termui/widgets"
	"gopkg.in/yaml.v2"
)

func trak() {

	var arg2 = flag.Args()[1]

	switch arg2 {
	case "alljobs":
		trakAllJobs()
	case "job":
		trakJob()
	default:
		log.Fatal("This subcommand does not exist!")
	}

}

// Makes POST Request and reads response
func trakAllJobs() {

	var location = flag.Args()[2]

	f, err := ioutil.ReadFile("./sharpcd.yml")
	var con config
	err = yaml.Unmarshal(f, &con)
	handle(err, "Failed to read and extract sharpcd.yml")

	url := con.Trak[location] + "/api/jobs/"
	secret := getSec()

	trakPayload := trakPostData{
		Secret:  secret,
		Version: sharpCDVersion,
	}

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

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	draw := func() {
		// Do Request
		resp, err := client.Do(req)
		handle(err, "Failed to do POST request"+url)
		defer resp.Body.Close()

		// Read Body and Status
		body, err := ioutil.ReadAll(resp.Body)
		handle(err, "Failed to read body of response"+url)

		var apiOutput response
		if err := json.Unmarshal(body, &apiOutput); err != nil {
			panic(err)
		}

		jobs := apiOutput.Jobs

		table1 := widgets.NewTable()
		tablePadding := 2

		table1.Rows = [][]string{
			[]string{"ID", "Name", "Type", "Status", "Error Message", "Reconnected"},
		}
		for _, job := range jobs {
			jobStr := []string{
				job.ID, job.Name, job.Type, job.Status, job.ErrMsg, strconv.FormatBool(job.Reconnect),
			}
			table1.Rows = append(table1.Rows, jobStr)
		}
		table1.ColumnWidths = generateColumns(table1.Rows, table1.Rows[0])
		width := sum(table1.ColumnWidths) + len(table1.Rows[0]) + 1

		p := widgets.NewParagraph()
		myFigure := figure.NewFigure("SharpCD Trak", "", true)
		figureRows := myFigure.Slicify()
		p.Text = ""
		figureSpace := 0

		if width < len(figureRows[0]) {
			width = len(figureRows[0]) + 4
			figureSpace = 0
		} else {
			figureSpace = (width-len(figureRows[0]))/2 - 2*tablePadding - 1
		}

		for _, row := range figureRows {
			p.Text += (strings.Repeat(" ", figureSpace) + row + "\n")
		}

		heightFigure := len(figureRows) + tablePadding
		p.SetRect(0, 0, width, heightFigure)
		ui.Render(p)

		deltaYTable := heightFigure
		heightTable := deltaYTable + len(table1.Rows)*2 + 1

		table1.TextStyle = ui.NewStyle(ui.ColorWhite)
		table1.SetRect(0, deltaYTable, width, heightTable)
		table1.TextAlignment = ui.AlignCenter
		table1.FillRow = true
		ui.Render(table1)

		deltaYClose := heightTable
		heightClose := deltaYClose + 3

		p = widgets.NewParagraph()
		p.Text = "Press Ctrl+C to Exit"
		p.SetRect(0, deltaYClose, width, heightClose)
		ui.Render(p)

	}

	draw()

	uiEvents := ui.PollEvents()
	ticker := time.NewTicker(time.Second).C
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

// For Tracking just a job
func trakJob() {

	var location = flag.Args()[2]
	var jobID = flag.Args()[3]

	f, err := ioutil.ReadFile("./sharpcd.yml")
	var con config
	err = yaml.Unmarshal(f, &con)
	handle(err, "Failed to read and extract sharpcd.yml")

	urlJob := con.Trak[location] + "/api/job/" + jobID
	urlLog := con.Trak[location] + "/api/logsfeed/" + jobID

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	var oldWidth = 0

	draw := func() {
		job := getAPIOutput(urlJob).Job

		table1 := widgets.NewTable()
		tablePadding := 2

		table1.Rows = [][]string{
			[]string{"ID", "Name", "Type", "Status", "Error Message", "Reconnected"},
		}
		jobStr := []string{
			job.ID, job.Name, job.Type, job.Status, job.ErrMsg, strconv.FormatBool(job.Reconnect),
		}
		table1.Rows = append(table1.Rows, jobStr)

		table1.ColumnWidths = generateColumns(table1.Rows, table1.Rows[0])
		width := sum(table1.ColumnWidths) + len(table1.Rows[0]) + 1

		if oldWidth > width {
			width = oldWidth
		} else {
			oldWidth = width
		}

		p := widgets.NewParagraph()
		myFigure := figure.NewFigure("SharpCD Trak", "", true)
		figureRows := myFigure.Slicify()
		p.Text = ""

		figureSpace := (width-len(figureRows[0]))/2 - 2*tablePadding - 1

		if figureSpace < 1 {
			figureSpace = 1
		}

		for _, row := range figureRows {
			p.Text += (strings.Repeat(" ", figureSpace) + row + "\n")
		}

		heightFigure := len(figureRows) + tablePadding
		p.SetRect(0, 0, width, heightFigure)
		ui.Render(p)

		deltaYTable := heightFigure
		heightTable := deltaYTable + len(table1.Rows)*2 + 1

		table1.TextStyle = ui.NewStyle(ui.ColorWhite)
		table1.SetRect(0, deltaYTable, width, heightTable)
		table1.TextAlignment = ui.AlignCenter
		table1.FillRow = true
		ui.Render(table1)

		logs := getAPIOutput(urlLog).Message

		deltaYLogs := heightTable
		heightLogs := deltaYLogs + 20

		p = widgets.NewParagraph()
		p.Text = logs
		p.SetRect(0, deltaYLogs, width, heightLogs)
		ui.Render(p)

		deltaYClose := heightLogs
		heightClose := deltaYClose + 3

		p = widgets.NewParagraph()
		p.Text = "Press Ctrl+C to Exit"
		p.SetRect(0, deltaYClose, width, heightClose)
		ui.Render(p)

	}

	draw()

	uiEvents := ui.PollEvents()
	ticker := time.NewTicker(time.Second).C
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

func sum(array []int) int {
	result := 0
	for _, v := range array {
		result += v
	}
	return result
}

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

func generateColumns(rows [][]string, columns []string) []int {

	columnWidths := []int{}

	for index := range columns {
		columnWidths = append(columnWidths, longest(rows, index))
	}

	return columnWidths
}

func getAPIOutput(url string) response {
	secret := getSec()

	trakPayload := trakPostData{
		Secret:  secret,
		Version: sharpCDVersion,
	}

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

	resp, err := client.Do(req)
	handle(err, "Failed to do POST request"+url)
	defer resp.Body.Close()

	// Read Body and Status
	body, err := ioutil.ReadAll(resp.Body)
	handle(err, "Failed to read body of response"+url)

	var apiOutput response
	if err := json.Unmarshal(body, &apiOutput); err != nil {
		panic(err)
	}

	return apiOutput
}
