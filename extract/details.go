/*
	The extract provides functions for downloading the package.json file and
	extracting the version of the required package from the body. It utilizes
	regular expressions for most of its functions and is directly called by the main
	package for storing the version
*/
package extract

import (
	"errors"
	"fmt"
	re "regexp"

	d "github.com/varuuntiwari/dyte-vit-2022-varuuntiwari/data"
)

func GetDetails(x d.RepoInput) (user, repo string, err error) {
	// Check if URL is valid
	if match, _ := re.MatchString("^(http|https)://github.com/[A-Za-z0-9-]{2,}/[A-Za-z0-9_-]+$", x.Link); !match {
		fmt.Println("no match")
		return
	}
	userRe := re.MustCompile(`github.com/(.*?)/[A-Za-z0-9_-]+$`)
	repoRe := re.MustCompile(`github.com/[A-Za-z0-9-_]+/(.*?)$`)
	err = nil

	if tmp := userRe.FindStringSubmatch(x.Link); tmp == nil {
		err = errors.New("cannot parse username")
		return
	} else {
		user = tmp[1]
	}
	if tmp := repoRe.FindStringSubmatch(x.Link); tmp == nil {
		err = errors.New("cannot parse repository name")
		return
	} else {
		repo = tmp[1]
	}
	return
}

// GetDependency calls functions within the package itself for
// retrieving the package.json for each GitHub repository after
// which it extracts the versions and returns whether the
// condition is satisfied or not.
func GetDependency(ver string, url string) (version string, fulfilled bool, err error) {
	resp, e := getPackage(url)
	if e != nil {
		err = errors.New("cannot get package.json")
		return
	}
	packName, packVersion := getSemanticVersion(ver)
	tmp := resp.(map[string]interface{})["dependencies"].(map[string]interface{})[packName]
	givenVersion := tmp.(string)[1:]
	return givenVersion, d.CompareVersions(packVersion, givenVersion), nil
}
