package main

const (
	statusNotPostMethod    = 606
	statusAcceptedTask     = 601
	statusFailedToReadBody = 607
	statusBodyNotJSON      = 608
	statusBannedURL        = 609
	statusIncorrectPass    = 610
	statusCommDoesNotExist = 612
)

type postData struct {
	ID         string            `json:"id"`
	Name       string            `json:"name"`
	Type       string            `json:"type"`
	GitURL     string            `json:"giturl"`
	Command    string            `json:"command"`
	Enviroment map[string]string `json:"enviroment"`
	Key        string            `json:"key"`
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
	Message string `json:"msg"`
}

type apiResponse struct {
	Procs []*taskProcess `json:"processes"`
}

type apiData struct {
	Key string `json:"key"`
}
