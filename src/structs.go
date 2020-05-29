package main

const (
	statusNotPostMethod    = 606
	statusAcceptedTask     = 601
	statusFailedToReadBody = 607
	statusBodyNotJSON      = 608
	statusBannedURL        = 609
	statusIncorrectPass    = 610
)

type postData struct {
	Name       string            `json:"name"`
	Type       string            `json:"type"`
	GitURL     string            `json:"giturl"`
	Command    string            `json:"command"`
	Enviroment map[string]string `json:"enviroment"`
	Key        string            `json:"key"`
}

type task struct {
	Name     string
	Type     string
	SharpURL string
	Envfile  string
	GitURL   string
	Command  string
}

type config struct {
	Version float32
	Tasks   map[string]task
}

type filter struct {
	Allowed []string
}

type response struct {
	Message string `json:"msg"`
}
