package cmd

import (
	"context"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/zkfmapf123/at-plan/usecase"
)

var rootCmd = &cobra.Command{
	Use:   "at-root",
	Short: "A CLI tool for managing your at-plan",
	Long:  `A CLI tool for managing your at-plan`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		githubToken, _ := cmd.Flags().GetString("github-token")
		githubRepoRef, _ := cmd.Flags().GetString("github-repo-ref")
		atlantisURL, _ := cmd.Flags().GetString("atlantis-url")
		atlantisToken, _ := cmd.Flags().GetString("atlantis-token")
		atlantisRepository, _ := cmd.Flags().GetString("atlantis-repository")
		atlantisConfig, _ := cmd.Flags().GetString("atlantis-config")
		slackBotToken, _ := cmd.Flags().GetString("slack-bot-token")
		slackChannel, _ := cmd.Flags().GetString("slack-channel")

		atRequest := usecase.AtlantisRequest{
			GithubToken:        githubToken,
			GithubRepoRef:      githubRepoRef,
			AtlantisURL:        atlantisURL,
			AtlantisToken:      atlantisToken,
			AtlantisRepository: atlantisRepository,
			AtlantisConfigFile: atlantisConfig,
			SlackBotToken:      slackBotToken,
			SlackChannel:       slackChannel,
		}

		ctx := context.WithValue(context.Background(), "atlantis", atRequest)
		cmd.SetContext(ctx)
	},
}

func init() {
	rootCmd.PersistentFlags().StringP("github-token", "g", "", "The Github token")
	rootCmd.PersistentFlags().StringP("github-repo-ref", "f", "", "The Github repository reference")
	rootCmd.PersistentFlags().StringP("atlantis-url", "u", "", "The Atlantis URL")
	rootCmd.PersistentFlags().StringP("atlantis-token", "t", "", "The Atlantis token")
	rootCmd.PersistentFlags().StringP("atlantis-repository", "r", "", "Atlantis Repository")
	rootCmd.PersistentFlags().StringP("atlantis-config", "c", "", "Atlantis Config File")
	rootCmd.PersistentFlags().StringP("slack-bot-token", "s", "", "Slack Bot Token")
	rootCmd.PersistentFlags().StringP("slack-channel", "l", "", "Slack Channel")
}

func Execute() {
	rootCmd.AddCommand(healthCmd)
	rootCmd.AddCommand(planCmd)
	rootCmd.AddCommand(notificationCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Println("Error executing command:", err)
		os.Exit(1)
	}
}
