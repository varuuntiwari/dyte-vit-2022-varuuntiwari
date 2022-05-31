package pulls

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/go-github/v45/github"
	ex "github.com/varuuntiwari/dyte-vit-2022-varuuntiwari/extract"
)

// cloneRepo takes a github Client program and forks the NodeJS project with an
// outdated dependency to the client account
func cloneRepo(user, repo string, cli *github.Client) (rep *github.Repository, err error) {
	repository := cli.Repositories
	rep, _, _ = repository.CreateFork(context.Background(), user, repo, nil)
	return
}

func updateRepoDependency(repo *github.Repository, cli *github.Client, version string) (err error) {
	body, e := ex.GetPackage(*repo.FullName)
	if e != nil {
		err = e
		return
	}
	packName, neededVersion := ex.GetSemanticVersion(version)
	body.(map[string]interface{})["dependencies"].(map[string]interface{})[packName] = neededVersion

	data, e := json.Marshal(body)
	if e != nil {
		err = e
		return
	}
	fmt.Println(string(data))

	u, _, _ := cli.Users.Get(context.Background(), "")
	cont, res, e := cli.Repositories.UpdateFile(context.Background(), *u.Login, *repo.Name, "package.json", &github.RepositoryContentFileOptions{
		Message: github.String("update dependency"),
		Content: data,
		SHA:     new(string),
		Branch:  github.String("depupdate"),
	})
	if e != nil {
		err = e
		return
	}
	fmt.Println(cont)
	fmt.Println(res)
	return
}
