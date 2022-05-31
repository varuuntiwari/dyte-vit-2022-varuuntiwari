package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"

	c "github.com/TreyBastian/colourize"
	d "github.com/varuuntiwari/dyte-vit-2022-varuuntiwari/data"
	ex "github.com/varuuntiwari/dyte-vit-2022-varuuntiwari/extract"
	"github.com/varuuntiwari/dyte-vit-2022-varuuntiwari/pulls"
)

var (
	csvFile       string
	update        bool
	neededVersion string
	cont          string
	authToken     string
	repos         []d.RepoInput
	results       []d.RepoOutput
)

func main() {
	// Command line parameters
	fmt.Println()
	flag.StringVar(&csvFile, "i", "", "Give CSV file as input")
	flag.BoolVar(&update, "update", false, "Update outdated dependencies")
	flag.StringVar(&authToken, "auth", "", "Enter authentication token")
	flag.StringVar(&neededVersion, "version", "", "Enter dependency and version to check")
	flag.Parse()

	if neededVersion == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	if csvFile == "" {
		flag.PrintDefaults()
		os.Exit(1)
	} else if match, _ := regexp.MatchString("^[a-z-]+@[0-9]+.[0-9]+.[0-9]+$", neededVersion); !match {
		fmt.Println(c.Colourize("[-] dependency not in valid format", c.Red))
		os.Exit(1)
	}

	if authToken == "" {
		fmt.Println(c.Colourize("auth token not given", c.Cyan))
		fmt.Println(c.Colourize("attempting to get token from ENV..", c.Cyan))
		envKey := "GITHUB_AUTH"
		if _, err := os.LookupEnv(envKey); !err {
			fmt.Println(c.Colourize("[-] token not found", c.Red))
			flag.PrintDefaults()
			os.Exit(1)
		} else if authToken = os.Getenv(envKey); authToken == "" {
			fmt.Println(c.Colourize("[-] token empty", c.Red))
			fmt.Println()
			flag.PrintDefaults()
			os.Exit(1)
		}
	}
	fmt.Println(c.Colourize("[+] auth token added", c.Green))

	// Confirm parameters given by user
	fmt.Printf("The input CSV file is %v\n", csvFile)
	fmt.Printf("Should the outdated dependencies be updated? %v\n", update)
	fmt.Printf("The dependency to be checked is: %v\n", neededVersion)

	fmt.Print("Continue?(y/n): ")
	cont = "y"
	fmt.Scanf("%v", &cont)

	if cont != "y" {
		fmt.Println(c.Colourize("[-] exiting..", c.Red))
		return
	}
	fmt.Println()

	// Opening CSV file and checking for availability
	f, err := os.Open(csvFile)
	if err != nil {
		fmt.Println(c.Colourize("[-] error opening file", c.Red))
		os.Exit(1)
	} else {
		fmt.Println(c.Colourize("[+] file readable", c.Green))
	}
	defer f.Close()

	// Test if file is readable into RepoInput struct
	reader := csv.NewReader(f)
	fmt.Println(c.Colourize("[+] reading file..", c.Green))
	for {
		repo, err := reader.Read()
		if err == io.EOF {
			fmt.Println(c.Colourize("[+] end of file", c.Green))
			break
		} else if err != nil {
			fmt.Println(c.Colourize("[-] error reading record", c.Green))
			break
		}

		repos = append(repos, d.RepoInput{
			Name: repo[0],
			Link: repo[1],
		})
	}

	// Display all repositories
	for _, repo := range repos {
		u, r, err := ex.GetDetails(repo)
		if err != nil {
			fmt.Println(c.Colourize("[-] "+err.Error(), c.Red))
		}
		url := u + "/" + r
		version, fulfilled, err := ex.GetDependency(neededVersion, url)
		if err != nil {
			fmt.Println(c.Colourize(err, c.Red))
			os.Exit(1)
		}
		prUrl := ""
		if !fulfilled {
			prUrl, err = pulls.CreatePR(authToken, u, r, neededVersion)
			if err != nil {
				fmt.Println(c.Colourize("[-]"+err.Error(), c.Red))
			}
		}
		results = append(results, d.RepoOutput{
			Name:      repo.Name,
			Link:      repo.Link,
			Version:   version,
			Satisfied: fulfilled,
			PRUpdate:  prUrl,
		})
	}
	for _, res := range results {
		fmt.Printf("%v\t%v\t%v\t%v\t%v\n", res.Name, res.Link, res.Version, res.Satisfied, res.PRUpdate)
	}
}
