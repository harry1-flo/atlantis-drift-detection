package usecase

type AtlantisRequestParams struct {
	BranchRef       string
	BranchName      string
	RepoName        string
	RepoCommitHash  string
	PRNum           string
	PRURL           string
	PRAuthor        string
	GHToken         string
	Owner           string
	ATCommand       string // validate, plan, apply
	SlackWebhookURL string
	RepoRelDir      string // 작업 파일 위치
	ChannelName     string
}

type PRParams struct {
	URL        string `json:"url"` // pr link
	ID         string `json:"id"`
	Number     string `json:"number"`
	State      string `json:"state"`
	Title      string `json:"title"`
	RepoRelDir string `json:"repo_rel_dir"` // 작업 파일 위치

	// PR Info
	ChangeFileCount string `json:"changed_files"`
	PRComments      string `json:"comments"`
	Commits         string `json:"commits"`

	// users
	Pusher     string
	PushCommit string

	// at
	Command string
}
