package client

import (
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

func NewGithubRequest(parms usecase.AtlantisRequestParams) *GithubParmas {
	return &GithubParmas{
		Request:    &parms,
		httpClient: utils.NewATHTTP(),
	}
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
