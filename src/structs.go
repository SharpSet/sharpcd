package main

import (
	"github.com/gorilla/websocket"
)

var sharpCDVersion = "3.6"

type statusCodes struct {
	NotPostMethod    int
	Accepted         int
	FailedToReadBody int
	BodyNotJSON      int
	BannedURL        int
	IncorrectSecret  int
	CommDoesNotExist int
	FailedLogRead    int
	WrongVersion     int
	JobDoesNotExist  int
}

var statCode = statusCodes{
	NotPostMethod:    606,
	Accepted:         607,
	FailedToReadBody: 608,
	BodyNotJSON:      612,
	IncorrectSecret:  613,
	FailedLogRead:    614,
	BannedURL:        616,
	CommDoesNotExist: 617,
	WrongVersion:     618,
	JobDoesNotExist:  619}

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

var allJobs map[string]*taskJob = make(map[string]*taskJob)

type jobStats struct {
	Running  string
	Errored  string
	Stopped  string
	Building string
	Stopping string
}

var jobStatus = jobStats{
	Running:  "running",
	Errored:  "errored",
	Stopped:  "stopped",
	Building: "building",
	Stopping: "stopping"}

type taskJob struct {
	ID         string            `json:"id"`
	Name       string            `json:"name"`
	Type       string            `json:"type"`
	Status     string            `json:"status"`
	ErrMsg     string            `json:"err_msg"`
	URL        string            `json:"url"`
	Enviroment map[string]string `json:"-"`
	Registry   string            `json:"registry"`
	Issue      string            `json:"issue"`
	Reconnect  bool
	Kill       bool
}

type postData struct {
	ID         string            `json:"id"`
	Name       string            `json:"name"`
	Type       string            `json:"type"`
	GitURL     string            `json:"giturl"`
	Command    string            `json:"command"`
	Enviroment map[string]string `json:"enviroment"`
	Secret     string            `json:"secret"`
	Compose    string            `json:"compose"`
	Registry   string            `json:"registry"`
	Version    string            `json:"version"`
}

type trakPostData struct {
	Secret  string `json:"secret"`
	Version string `json:"version"`
}

type task struct {
	ID       string
	Name     string
	Type     string
	SharpURL string
	Envfile  string
	GitURL   string
	Command  string
	Compose  string
	Registry string
	Depends  []string
}

type config struct {
	Version float32
	Tasks   map[string]task
	Trak    map[string]string
}

type filter struct {
	Allowed []string
	Token   string
}

type response struct {
	Jobs    []*taskJob `json:"jobs"`
	Job     *taskJob   `json:"job"`
	Message string     `json:"msg"`
}

type allDirs struct {
	Root    string
	Private string
	Logs    string
	Docker  string
}

var folder = allDirs{
	Root:    getDir() + "/sharpcd-data/",
	Private: getDir() + "/sharpcd-data/private/",
	Docker:  getDir() + "/sharpcd-data/docker/"}
