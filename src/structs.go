package main

type statusCodes struct {
	NotPostMethod    int
	Accepted         int
	FailedToReadBody int
	BodyNotJSON      int
	BannedURL        int
	IncorrectSecret  int
	CommDoesNotExist int
	FailedLogRead    int
}

var statCode = statusCodes{
	NotPostMethod:    606,
	Accepted:         607,
	FailedToReadBody: 608,
	BodyNotJSON:      612,
	IncorrectSecret:  613,
	FailedLogRead:    614,
	BannedURL:        616,
	CommDoesNotExist: 617}

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
	ID     string `json:"id"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	Status string `json:"status"`
	ErrMsg string `json:"err_msg"`
	URL    string `json:"url"`
}

type allTaskJobs struct {
	List []*taskJob
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
}

type config struct {
	Version float32
	Tasks   map[string]task
}

type filter struct {
	Allowed []string
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
	Logs:    getDir() + "/sharpcd-data/logs/",
	Docker:  getDir() + "/sharpcd-data/docker/"}
