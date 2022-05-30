package extract

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

// Download the package.json file from given GitHub link
// and return the body of the response received.
func getPackage(u string) (res interface{}, err error) {
	branch := "main"
	rootURL := "https://raw.githubusercontent.com/" + u + "/" + branch + "/package.json"

	resp, e := http.Get(rootURL)
	if e != nil || resp.StatusCode != 200 {
		err = errors.New("cannot access url of repository")
		return
	}
	defer resp.Body.Close()
	e = json.NewDecoder(resp.Body).Decode(&res)
	if e != nil {
		err = errors.New("cannot decode file to json")
		return
	}
	return
}

// Get the package name and version from the format
// it is stored inside package.json.
func getSemanticVersion(ver string) (string, string) {
	s := strings.Split(ver, "@")
	return s[0], s[1]
}
