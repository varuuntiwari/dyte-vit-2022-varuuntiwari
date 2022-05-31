package pulls

import (
	"context"
	"fmt"
	"time"

	c "github.com/TreyBastian/colourize"
	"github.com/google/go-github/v45/github"
	"golang.org/x/oauth2"
)

// CreatePR creates a pull request for a given repository with authentication token
func CreatePR(token, user, repo, ver string) (prUrl string, err error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	u, _, e := client.Users.Get(ctx, "")
	if e != nil {
		fmt.Println(c.Colourize("[-] cannot access username", c.Red))
	} else {
		fmt.Printf(c.Colourize("Username of access token: %v\n", c.Green), *u.Login)
	}

	// call the cloneRepo to fork repository and make changes
	forked, e := cloneRepo(user, repo, client)
	if e != nil {
		err = e
		return
	}
	// Wait for fork to be updated
	fmt.Println(c.Colourize("waiting for fork to be updated in account..", c.Cyan))
	time.Sleep(5 * time.Second)

	// Update the package.json to include package version required
	e = updateRepoDependency(forked, client, ver)
	if e != nil {
		err = e
		return
	}

	newPR := &github.NewPullRequest{
		Title:               github.String("PR for dependency update"),
		Head:                github.String(*u.Name + ":main"),
		Base:                github.String("main"),
		Body:                github.String("update dependency to " + ver),
		MaintainerCanModify: github.Bool(true),
	}
	pr, _, e := client.PullRequests.Create(context.Background(), user, repo, newPR)
	if e != nil {
		err = e
		return
	}
	prUrl = pr.GetHTMLURL()
	return
}
