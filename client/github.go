package client

import (
	"errors"
	"fmt"
	"log"

	"github.com/zkfmapf123/at-plan/usecase"
	"github.com/zkfmapf123/at-plan/utils"
	"github.com/zkfmapf123/donggo"
)

type GithubParmas struct {
	Request *usecase.AtlantisRequestParams

	httpClient *utils.ATHTTP
}

func NewGithubRequest(parms usecase.AtlantisRequestParams) (*GithubParmas, error) {

	params := &GithubParmas{
		Request:    &parms,
		httpClient: utils.NewATHTTP(),
	}

	err := githubParamValidate(*params)

	return params, err
}

func githubParamValidate(params GithubParmas) error {

	if params.Request.GHToken == "" {
		return errors.New("github token is empty")
	}
	if params.Request.PRNum == "" {
		return errors.New("pr number is empty")
	}
	if params.Request.PRURL == "" {
		return errors.New("pr url is empty")
	}
	if params.Request.PRAuthor == "" {
		return errors.New("pr author is empty")
	}
	if params.Request.ATCommand == "" {
		return errors.New("at command is empty")
	}
	if params.Request.RepoRelDir == "" {
		return errors.New("repo rel dir is empty")
	}
	if params.Request.SlackBotToken == "" {
		return errors.New("slack bot token is empty")
	}
	if params.Request.SlackChannel == "" {
		return errors.New("slack channel is empty")
	}
	if params.Request.Outputs == "" {
		return errors.New("outputs is empty")
	}
	if params.Request.Owner == "" {
		return errors.New("owner is empty")
	}
	if params.Request.RepoName == "" {
		return errors.New("repo name is empty")
	}
	if params.Request.RepoCommitHash == "" {
		return errors.New("repo commit hash is empty")
	}
	if params.Request.BranchRef == "" {
		return errors.New("branch ref is empty")
	}
	if params.Request.BranchName == "" {
		return errors.New("branch name is empty")
	}

	return nil
}

// 현재 새로운 PR 인지 여부
// PR Comments가 몇개인지로 판단...
func (g GithubParmas) IsNewPR() (usecase.PRParams, bool) {

	resp, err := g.httpClient.Comm(
		utils.HTTPParams{
			Url:    fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls/%s", g.Request.Owner, g.Request.RepoName, g.Request.PRNum),
			Method: "GET",
			Headers: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", g.Request.GHToken),
				"Accept":        "application/vnd.github+json",
			},
		},
	)

	if err != nil {
		panic(err)
	}

	result := donggo.JsonParse[usecase.PRParams](resp)

	log.Printf("Git PR Number : %d Status : %s Comments Count : %d", result.Number, result.State, result.PRComments)

	return result, result.State == "open" && result.Commits == 0
}
