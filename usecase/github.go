package usecase

import "fmt"

var (
	VALIDATE = "validate"
	PLAN     = "plan"
	APPLY    = "apply"

	ERROR  = "error"
	FAILED = "failed"
)

type AtlantisRequestParams struct {
	BranchRef      string
	BranchName     string
	RepoName       string
	RepoCommitHash string
	PRNum          string
	PRURL          string
	PRAuthor       string
	GHToken        string
	Owner          string
	ATCommand      string // validate, plan, apply
	RepoRelDir     string // 작업 파일 위치
	SlackBotToken  string
	SlackChannel   string
	Outputs        string
}

type PRParams struct {
	URL        string `json:"url"` // pr link
	ID         string `json:"id"`
	Number     int    `json:"number"`
	State      string `json:"state"` // Open, Closed
	Title      string `json:"title"`
	RepoRelDir string `json:"repo_rel_dir"` // 작업 파일 위치

	// PR Info
	ChangeFileCount int `json:"changed_files"`
	PRComments      int `json:"comments"`
	Commits         int `json:"commits"`

	// users
	Pusher     string
	PushCommit string

	// at-slack
	Command       string
	SlackBotToken string
	SlackChannel  string

	Outputs string
}

func NewPR(command string, PRNum int) string {
	return fmt.Sprintf("New Pull Request Created [%s] #%d", command, PRNum)
}

func Validate(command string, PRNum int) string {
	return fmt.Sprintf("Validate Pull Request [#%d]", PRNum)
}

func Plan(command string, PRNum int) string {
	return fmt.Sprintf("Atlantis Plan [#%d]", PRNum)
}

func Apply(command string, PRNum int) string {
	return fmt.Sprintf("Atlantis Apply [#%d]", PRNum)
}
