package client

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/stretchr/testify/assert/yaml"
	"github.com/zkfmapf123/at-plan/usecase"
	"github.com/zkfmapf123/at-plan/utils"
	"github.com/zkfmapf123/donggo"
)

/*
@reference
https://www.runatlantis.io/docs/api-endpoints
*/
var (
	API_PLAN   = "/api/plan"
	API_HEALTH = "/healthz"
)

type AtlantisParams struct {
	Request *usecase.AtlantisRequest

	atlantisConfigParmas usecase.AtlantisConfigParams
	httpClient           *utils.ATHTTP
}

func NewAtlantisRequest(atParams usecase.AtlantisRequest) AtlantisParams {
	return AtlantisParams{
		Request:    &atParams,
		httpClient: utils.NewATHTTP(),
	}
}

func (a *AtlantisParams) ValidURL() error {

	resp, err := a.httpClient.Comm(
		utils.HTTPParams{
			Url:    a.Request.AtlantisURL + API_HEALTH,
			Method: "GET",
		},
	)

	if err != nil {
		return err
	}

	res := donggo.JsonParse[usecase.APIHealthResponse](resp)
	if res.Status != "ok" {
		return errors.New("atlantis is not running")
	}

	return nil
}

func (a *AtlantisParams) ValidRepository() error {

	gitAttr := strings.Split(a.Request.AtlantisRepository, "/")
	owner, repository := gitAttr[0], gitAttr[1]

	resp, err := a.httpClient.Comm(
		utils.HTTPParams{
			Url:    fmt.Sprintf("https://api.github.com/repos/%s/%s", owner, repository),
			Method: "GET",
		},
	)

	if err != nil {
		return err
	}

	res := donggo.JsonParse[any](resp)
	if res == nil {
		return errors.New("repository not found")
	}

	return nil
}

func (a *AtlantisParams) ValidConfigFile() error {

	atFilePath := a.Request.AtlantisConfigFile

	// 상대 경로면 절대 경로로 변환
	if !filepath.IsAbs(atFilePath) {
		pwd, _ := utils.GetPwd()
		atFilePath = filepath.Join(pwd, atFilePath)
	}

	_, err := os.Stat(atFilePath)
	if os.IsNotExist(err) {
		return errors.New("atlantis config file not found")
	}

	// rewrite
	a.Request.AtlantisConfigFile = atFilePath
	return nil
}

func (a *AtlantisParams) SetConfigParmas() (usecase.AtlantisConfigParams, error) {

	var params usecase.AtlantisConfigParams

	if err := a.ValidConfigFile(); err != nil {
		return usecase.AtlantisConfigParams{}, err
	}

	b, err := os.ReadFile(a.Request.AtlantisConfigFile)
	if err != nil {
		return usecase.AtlantisConfigParams{}, err
	}

	if err := yaml.Unmarshal(b, &params); err != nil {
		return usecase.AtlantisConfigParams{}, err
	}

	a.atlantisConfigParmas = params
	return params, nil
}

func (a *AtlantisParams) Plan() {

	// gitAttr := strings.Split(a.Request.AtlantisRepository, "/")
	// // owner := gitAttr[0]
	// repository := gitAttr[1]

	for _, project := range a.atlantisConfigParmas.Projects {

		resp, err := a.httpClient.Comm(
			utils.HTTPParams{
				Url:    fmt.Sprintf("%s%s", a.Request.AtlantisURL, API_PLAN),
				Method: "POST",
				Headers: map[string]string{
					"X-Atlantis-Token": a.Request.AtlantisToken,
					"Content-Type":     "application/json",
				},
				Body: map[string]any{
					"Repository": a.Request.AtlantisRepository,
					"Ref":        a.Request.GithubRepoRef,
					"Type":       "Github",
					"Paths": []usecase.APIPlanBodyPaths{
						{
							ProjectName: project.Name,
							Directory:   project.Dir,
							Workspace:   "default",
						},
					},
				},
			},
		)

		if err != nil {
			panic(err)
		}

		fmt.Println("error : ", err)
		fmt.Println(string(resp))
	}
}
