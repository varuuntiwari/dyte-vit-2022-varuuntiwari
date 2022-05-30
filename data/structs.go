/*
	This package contains the data structs for maintaining the structure of data
	taken from input and formatting the output properly.
*/
package data

import (
	"strings"
)

// Struct to store information from CSV file
type RepoInput struct {
	Name string
	Link string
}

// Struct to store information for printing
// as output
type RepoOutput struct {
	Name      string
	Link      string
	Version   string
	Satisfied bool
}

// Compare versions of required dependency version
// and version given in package.json
func CompareVersions(needed, given string) bool {
	neededVersion := strings.Split(needed, ".")
	givenVersion := strings.Split(given, ".")

	for i := 0; i < 3; i++ {
		if givenVersion[i] > neededVersion[i] {
			return true
		} else if givenVersion[i] == neededVersion[i] {
			continue
		} else if givenVersion[i] < neededVersion[i] {
			return false
		}
	}
	return true
}
