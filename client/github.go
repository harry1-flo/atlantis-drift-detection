package client

import (
	"errors"
	"fmt"
	"strings"

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

	if params.Request.BaseRepoName == "" {
		return errors.New("base repo name is required")
	}

	if params.Request.BaseRepoOwner == "" {
		return errors.New("base repo owner is required")
	}

	if params.Request.HeadCommit == "" {
		return errors.New("head commit is required")
	}

	if params.Request.PullURL == "" {
		return errors.New("pull url is required")
	}

	if params.Request.PullAuthor == "" {
		return errors.New("pull author is required")
	}

	if params.Request.Dir == "" {
		return errors.New("dir is required")
	}

	if params.Request.UserName == "" {
		return errors.New("user name is required")
	}

	if params.Request.CommandName == "" {
		return errors.New("command name is required")
	}

	if params.Request.GHToken == "" {
		return errors.New("github token is required")
	}

	if params.Request.SlackBotToken == "" {
		return errors.New("slack bot token is required")
	}

	if params.Request.SlackChannel == "" {
		return errors.New("slack channel is required")
	}

	return nil
}

// 현재 새로운 PR 인지 여부
// PR Comments가 몇개인지로 판단...
func (g GithubParmas) IsNewPR() (usecase.PRParams, bool) {

	fmt.Println(">>>", g.Request.BaseRepoOwner, g.Request.BaseRepoName, g.Request.PullURL)

	_PRNum := strings.Split(g.Request.PullURL, "/")
	PRNum := _PRNum[len(_PRNum)-1]

	resp, err := g.httpClient.Comm(
		utils.HTTPParams{
			Url:    fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls/%s", g.Request.BaseRepoOwner, g.Request.BaseRepoName, PRNum),
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
	return result, result.State == "open" && result.Commits == 0
}

func (g GithubParmas) GetCommentsLastPR(params usecase.PRParams) (status string, shortMessage string) {

	resp, err := g.httpClient.Comm(
		utils.HTTPParams{
			Url: fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls/%d/comments",
				g.Request.BaseRepoOwner, g.Request.BaseRepoName, params.Number),
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

	results := donggo.JsonParse[[]usecase.PRComments](resp)

	// PR
	if len(results) <= 1 {
		return "success", ""
	}

	// Get Last Comment
	lastResults := results[len(results)-1]

	// Set Status
	status = "success"
	if strings.Contains(lastResults.Body, "Error") {
		status = "failed"
	}

	// Extract Error message and Plan Summary
	var planSummary string

	// Error
	if strings.Contains(lastResults.Body, "Error") {
		return status, lastResults.Body
	}

	// Plan / Apply Summary
	for str := range strings.SplitSeq(shortMessage, "\n") {
		trimmedStr := strings.TrimSpace(str)

		if strings.Contains(trimmedStr, "projects") && strings.Contains(trimmedStr, "with changes") {
			planSummary = trimmedStr
		}
	}

	return status, planSummary
}
