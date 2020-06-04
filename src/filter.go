package main

import (
	"flag"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

var filterLoc = folder.Root + "filter.yml"

// Adds a url to filter.yml
func addFilter() {
	url := flag.Args()[1]

	// Reads filter.yml
	f, err := readFilter()
	handle(err, "Failed to read filter.yml")

	// Adds in new allowed url and saves
	f.Allowed = append(f.Allowed, url)
	err = saveFilter(f)
	handle(err, "Failed to save filter.yml")
}

// Removes a url to filter.yml
func removeFilter() {
	targetURL := flag.Args()[1]

	f, err := readFilter()
	handle(err, "Failed to read filter.yml")

	// Removes banned URL and saves
	var newAllowed []string
	for _, url := range f.Allowed {
		if url != targetURL {
			newAllowed = append(newAllowed, url)
		}
	}

	f.Allowed = newAllowed
	err = saveFilter(f)
	handle(err, "Failed to save filter.yml")
}

// Reads filter.yml
func readFilter() (filter, error) {
	var f filter
	var file []byte
	var err error

	// If file does not exist
	if !fileExists(filterLoc) {

		// Create it
		emptyfilter := filter{}
		saveFilter(emptyfilter)
	}

	// Read file
	file, err = ioutil.ReadFile(filterLoc)
	if err != nil {
		return f, err
	}

	// Load file into filter struct
	err = yaml.Unmarshal(file, &f)
	if err != nil {
		return f, err
	}

	return f, nil
}

// Save Filter
func saveFilter(f filter) error {

	// Marshal filter struct and save
	yml, err := yaml.Marshal(f)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(folder.Root+"/filter.yml", yml, 0644)
	if err != nil {
		return err
	}

	return nil
}

// fileExists checks if a file exists
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
