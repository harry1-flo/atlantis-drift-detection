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

		atReqParams.BaseRepoName, _ = cmd.Flags().GetString("base-repo-name")
		atReqParams.BaseRepoOwner, _ = cmd.Flags().GetString("base-repo-owner")
		atReqParams.HeadCommit, _ = cmd.Flags().GetString("head-commit")
		atReqParams.PullURL, _ = cmd.Flags().GetString("pull-url")
		atReqParams.PullAuthor, _ = cmd.Flags().GetString("pull-author")
		atReqParams.Dir, _ = cmd.Flags().GetString("dir")
		atReqParams.UserName, _ = cmd.Flags().GetString("user-name")
		atReqParams.CommandName, _ = cmd.Flags().GetString("command-name")
		atReqParams.GHToken, _ = cmd.Flags().GetString("gh-token")
		atReqParams.SlackBotToken, _ = cmd.Flags().GetString("slack-bot-token")
		atReqParams.SlackChannel, _ = cmd.Flags().GetString("slack-channel")
		atReqParams.CommandHasErrors, _ = cmd.Flags().GetBool("command-has-errors")

		gc, err := client.NewGithubRequest(atReqParams)
		if err != nil {
			log.Fatalf("github request error: %s", err)
		}

		prParams, isNewPR := gc.IsNewPR()
		// status, shortMessage := gc.GetCommentsLastPR(prParams)

		prParams.Command = atReqParams.CommandName
		prParams.SlackBotToken = atReqParams.SlackBotToken
		prParams.SlackChannel = atReqParams.SlackChannel
		prParams.URL = atReqParams.PullURL
		prParams.PushCommit = atReqParams.HeadCommit
		prParams.Pusher = atReqParams.PullAuthor

		log.Println("isError : ", atReqParams.CommandHasErrors)

		status := "success"
		if !atReqParams.CommandHasErrors {
			status = "failed"
		}

		// PR 처음인 경우
		if isNewPR {
			err := utils.SendSlackAtlantisNoti(prParams, status)
			if err != nil {
				log.Fatalf("slack send error: %s", err)
			}
			return
		}

		// /*
		// 	init / validate 실패
		// 	plan 실패
		// 	apply 성공
		// 	apply 실패
		// */
		if status == "failed" {
			if err := utils.SendSlackAtlantisNoti(prParams, status); err != nil {
				log.Fatalf("slack send error: %s", err)
			}

			return
		}

		return
	},
}

func init() {
	notificationCmd.Flags().String("base-repo-name", "", "The base repository name")
	notificationCmd.Flags().String("base-repo-owner", "", "The base repository owner")
	notificationCmd.Flags().String("head-commit", "", "The head commit hash")
	notificationCmd.Flags().String("pull-url", "", "The pull request URL")
	notificationCmd.Flags().String("pull-author", "", "The pull request author")
	notificationCmd.Flags().String("dir", "", "The directory")
	notificationCmd.Flags().String("user-name", "", "The user name")
	notificationCmd.Flags().String("command-name", "", "The command name")
	notificationCmd.Flags().Bool("command-has-errors", false, "Whether the command has errors")

	notificationCmd.Flags().String("slack-bot-token", "", "The Slack bot token")
	notificationCmd.Flags().String("slack-channel", "", "The Slack channel")
	notificationCmd.Flags().String("gh-token", "", "The Github token")
}
