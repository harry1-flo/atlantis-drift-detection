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
var notificationCmd = &cobra.Command{
	Use:   "notification",
	Short: "A CLI tool for managing your notification",
	Long:  `A CLI tool for managing your notification`,
	Run: func(cmd *cobra.Command, args []string) {

		var atReqParams usecase.AtlantisRequestParams
		atReqParams.BranchRef, _ = cmd.Flags().GetString("at-branch-ref")
		atReqParams.BranchName, _ = cmd.Flags().GetString("at-branch-name")
		atReqParams.RepoName, _ = cmd.Flags().GetString("at-repo-name")
		atReqParams.RepoCommitHash, _ = cmd.Flags().GetString("at-commit-hash")
		atReqParams.PRNum, _ = cmd.Flags().GetString("at-pr-num")
		atReqParams.PRURL, _ = cmd.Flags().GetString("at-pr-url")
		atReqParams.PRAuthor, _ = cmd.Flags().GetString("at-pr-author")
		atReqParams.GHToken, _ = cmd.Flags().GetString("at-gh-token")
		atReqParams.ATCommand, _ = cmd.Flags().GetString("at-command")
		atReqParams.Owner, _ = cmd.Flags().GetString("at-owner")
		atReqParams.RepoRelDir, _ = cmd.Flags().GetString("at-repo-rel-dir")
		atReqParams.SlackBotToken, _ = cmd.Flags().GetString("at-slack-bottoken")
		atReqParams.SlackChannel, _ = cmd.Flags().GetString("at-slack-channel")
		atReqParams.Outputs, _ = cmd.Flags().GetString("at-outputs")

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

		log.Println("atReqParams.Outputs : ", atReqParams.Outputs)
		log.Println("parsing Terraform Outputs : ", msg)
		log.Println(">>>>>>>>>>>>>> isNewPR : ", isNewPR, "status : ", status, "command : ", prParams.Command)

		// PR 처음인 경우
		if isNewPR {
			utils.SendSlackAtlantisNoti(prParams, status)
			return
		}

		/*
			init / validate 실패
			plan 실패
			apply 성공
			apply 실패
		*/
		if status == "failed" && (prParams.Command == usecase.PLAN || prParams.Command == usecase.APPLY) {
			utils.SendSlackAtlantisNoti(prParams, status)
			return
		}

		return
	},
}

func init() {
	notificationCmd.Flags().String("at-branch-ref", "", "The Atlantis branch reference")
	notificationCmd.Flags().String("at-branch-name", "", "The Atlantis branch name")
	notificationCmd.Flags().String("at-repo-name", "", "The Atlantis repository name")
	notificationCmd.Flags().String("at-commit-hash", "", "The Atlantis commit hash")
	notificationCmd.Flags().String("at-pr-num", "", "The Atlantis PR number")
	notificationCmd.Flags().String("at-pr-url", "", "The Atlantis PR URL")
	notificationCmd.Flags().String("at-pr-author", "", "The Atlantis PR author")
	notificationCmd.Flags().String("at-gh-token", "", "The Github token")
	notificationCmd.Flags().String("at-command", "", "The Atlantis command")
	notificationCmd.Flags().String("at-owner", "", "The Atlantis owner")
	notificationCmd.Flags().String("at-repo-rel-dir", "", "The Atlantis repository relative directory")
	notificationCmd.Flags().String("at-slack-bottoken", "", "The Atlantis slack webhook URL")
	notificationCmd.Flags().String("at-slack-channel", "", "The Atlantis slack channel")
	notificationCmd.Flags().String("at-outputs", "", "The Atlantis outputs")

}
