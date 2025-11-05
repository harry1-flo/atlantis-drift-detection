package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/zkfmapf123/at-plan/client"
	"github.com/zkfmapf123/at-plan/usecase"
	"github.com/zkfmapf123/at-plan/utils"
)

/*
요구사항 정의
PR 생성 시

	-> Init 실패 / 성공 여부 확인
	-> validate 여부 확인 (계속 확인되어야 함)
	-> Apply 시, 실패 성공 여부 확인

슬랙메시지 시, 아래 스레드로 달릴 수 있는지?
*/
var githubCmd = &cobra.Command{
	Use:   "github",
	Short: "A CLI tool for managing your github",
	Long:  `A CLI tool for managing your github`,
	Run: func(cmd *cobra.Command, args []string) {

		var atReqParams usecase.AtlantisRequestParams
		atReqParams.BranchRef = cmd.PersistentFlags().Lookup("at-branch-ref").Value.String()
		atReqParams.BranchName = cmd.PersistentFlags().Lookup("at-branch-name").Value.String()
		atReqParams.RepoName = cmd.PersistentFlags().Lookup("at-repo-name").Value.String()
		atReqParams.RepoCommitHash = cmd.PersistentFlags().Lookup("at-commit-hash").Value.String()
		atReqParams.PRNum = cmd.PersistentFlags().Lookup("at-pr-num").Value.String()
		atReqParams.PRURL = cmd.PersistentFlags().Lookup("at-pr-url").Value.String()
		atReqParams.PRAuthor = cmd.PersistentFlags().Lookup("at-pr-author").Value.String()
		atReqParams.GHToken = cmd.PersistentFlags().Lookup("at-gh-token").Value.String()
		atReqParams.ATCommand = cmd.PersistentFlags().Lookup("at-command").Value.String()
		atReqParams.Owner = cmd.PersistentFlags().Lookup("at-owner").Value.String()
		atReqParams.RepoRelDir = cmd.PersistentFlags().Lookup("at-repo-rel-dir").Value.String()
		atReqParams.SlackBotToken = cmd.PersistentFlags().Lookup("at-slack-bottoken").Value.String()
		atReqParams.SlackChannel = cmd.PersistentFlags().Lookup("at-slack-channel").Value.String()
		atReqParams.Outputs = cmd.PersistentFlags().Lookup("at-outputs").Value.String()

		gc, err := client.NewGithubRequest(atReqParams)
		if err != nil {
			log.Fatalf("github request error: %s", err)
		}

		prParams, isNewPR := gc.IsNewPR()
		prParams.Pusher = atReqParams.PRAuthor
		prParams.PushCommit = atReqParams.RepoCommitHash
		prParams.Command = atReqParams.ATCommand
		prParams.SlackBotToken = atReqParams.SlackBotToken
		prParams.SlackChannel = atReqParams.SlackChannel
		prParams.Outputs = atReqParams.Outputs
		prParams.RepoRelDir = atReqParams.RepoRelDir

		status, msg := utils.LinseToParseLastMesasge(atReqParams.Outputs)
		prParams.Outputs = msg

		// PR 처음인 경우
		if isNewPR {

			utils.SendSlackAtlantisNoti(prParams, status)
		}

		// init / validate 실패
		if atReqParams.ATCommand == usecase.VALIDATE {

		}

		// plan 실패
		if atReqParams.ATCommand == usecase.PLAN {

		}

		// apply 성공
		if atReqParams.ATCommand == usecase.APPLY {

		}

		// apply 실패
		if atReqParams.ATCommand == usecase.APPLY {

		}
	},
}

func init() {
	githubCmd.PersistentFlags().String("at-branch-ref", "", "The Atlantis branch reference")
	githubCmd.PersistentFlags().String("at-branch-name", "", "The Atlantis branch name")
	githubCmd.PersistentFlags().String("at-repo-name", "", "The Atlantis repository name")
	githubCmd.PersistentFlags().String("at-commit-hash", "", "The Atlantis commit hash")
	githubCmd.PersistentFlags().String("at-pr-num", "", "The Atlantis PR number")
	githubCmd.PersistentFlags().String("at-pr-url", "", "The Atlantis PR URL")
	githubCmd.PersistentFlags().String("at-pr-author", "", "The Atlantis PR author")
	githubCmd.PersistentFlags().String("at-gh-token", "", "The Github token")
	githubCmd.PersistentFlags().String("at-command", "", "The Atlantis command")
	githubCmd.PersistentFlags().String("at-owner", "", "The Atlantis owner")
	githubCmd.PersistentFlags().String("at-repo-rel-dir", "", "The Atlantis repository relative directory")
	githubCmd.PersistentFlags().String("at-slack-bottoken", "", "The Atlantis slack webhook URL")
	githubCmd.PersistentFlags().String("at-slack-channel", "", "The Atlantis slack channel")
	githubCmd.PersistentFlags().String("at-outputs", "", "The Atlantis outputs")

}
