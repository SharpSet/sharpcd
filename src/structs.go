package main

type postData struct {
	Name       string            `json:"name"`
	Type       string            `json:"type"`
	GitURL     string            `json:"giturl"`
	Command    string            `json:"command"`
	Enviroment map[string]string `json:"enviroment"`
	Key        string            `json:"key"`
}

type task struct {
	Name string
	Type string
	SharpURL string
	Envfile string
	GitURL string
	Command string
}

type config struct {
    Version float32
    Tasks map[string]task
}


