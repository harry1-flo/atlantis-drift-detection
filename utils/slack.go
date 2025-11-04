package utils

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/slack-go/slack"
	"github.com/zkfmapf123/at-plan/usecase"
)

func SendSlack(
	slackBotToken string,
	slackChannel string,
	tfOutputs map[string]string,
) error {

	if slackBotToken == "" || slackChannel == "" {
		return errors.New("slack bot token or channel is empty")
	}

	if len(tfOutputs) == 0 {
		return errors.New("no terraform outputs provided")
	}

	planOutputContent := createContent(tfOutputs)
	fmt.Printf("Content size: %d bytes\n", len(planOutputContent))
	if len(planOutputContent) == 0 {
		return errors.New("generated content is empty")
	}

	//////////////////////////////////////////// Create ////////////////////////////////////////////
	api := slack.New(slackBotToken)
	reader := bytes.NewReader([]byte(planOutputContent))

	placeHolder := fmt.Sprintf("Drift detected in %d project(s)\n", countNonEmptyOutputs(tfOutputs))
	for dir, output := range tfOutputs {
		add, change, destroy := linesToPlanOutput(output)

		if add == "" && change == "" && destroy == "" {
			placeHolder += fmt.Sprintf("📁 Project : %s - No Changes\n", dir)
		} else {
			placeHolder += fmt.Sprintf("📁 Project : %s - Add: %s, Change: %s, Destroy: %s\n", dir, add, change, destroy)

		}
	}

	params := slack.UploadFileV2Parameters{
		Title:          "Drift Detection Report",
		InitialComment: placeHolder,
		FileSize:       len(planOutputContent), // 파일 크기 명시
		Reader:         reader,                 // Content 대신 Reader
		Filename:       "drift-report.txt",
		Channel:        slackChannel,
	}

	_, err := api.UploadFileV2(params)
	if err != nil {
		return fmt.Errorf("slack upload failed: %w", err)
	}

	fmt.Println("✅ Slack file uploaded successfully!")
	return nil
}

func createContent(tfOutputs map[string]string) string {
	var b strings.Builder

	b.WriteString("╔════════════════════════════════════════════════════════════╗\n")
	b.WriteString("║          ATLANTIS DRIFT DETECTION REPORT                   ║\n")
	b.WriteString("╚════════════════════════════════════════════════════════════╝\n\n")

	projectCount := 0
	for dir, output := range tfOutputs {
		if output == "" {
			continue
		}

		projectCount++
		b.WriteString(fmt.Sprintf("\n📁 Project %d: %s\n", projectCount, dir))
		b.WriteString("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n")
		b.WriteString(output)
		b.WriteString("\n\n")
	}

	b.WriteString("╔════════════════════════════════════════════════════════════╗\n")
	b.WriteString(fmt.Sprintf("║  Total Projects with Drift: %-31d║\n", projectCount))
	b.WriteString("╚════════════════════════════════════════════════════════════╝\n")

	return b.String()
}

func countNonEmptyOutputs(tfOutputs map[string]string) int {
	count := 0
	for _, output := range tfOutputs {
		if output != "" {
			count++
		}
	}
	return count
}

var (
	NEW      = "new"
	VALIDATE = "validate"
	PLAN     = "plan"
	APPLY    = "apply"
	ERROR    = "error"
	FAILED   = "failed"
)

// Slack Template Atlantis
func SendSlackAtlantisNoti(params usecase.PRParams, status string) error {

	if params.SlackBotToken == "" {
		return errors.New("slack bot token is empty")
	}

	api := slack.New(params.SlackBotToken)

	// Status에 따라 색상과 이모지 결정
	color := getStatusColor(status)
	emoji := getStatusEmoji(status)

	var title string
	switch status {

	case NEW:
		title = usecase.NewPR(params.Command, params.Number)
	case VALIDATE:
		title = usecase.Validate(params.Command, params.Number)
	case PLAN:
		title = usecase.Plan(params.Command, params.Number)
	case APPLY:
		title = usecase.Apply(params.Command, params.Number)
	}

	attachment := slack.Attachment{
		Color: color,
		Title: title,
		Fields: []slack.AttachmentField{
			{
				Title: "PR Link",
				Value: params.URL,
				Short: false,
			},
			{
				Title: "Status",
				Value: fmt.Sprintf("%s Terraform %s %s", emoji, getCommandTitle(params.Command), params.State),
				Short: false,
			},
			{
				Title: "Pusher",
				Value: params.Pusher,
				Short: true,
			},
			{
				Title: "Commit",
				Value: params.PushCommit[:7], // short commit hash
				Short: true,
			},
			{
				Title: "Project",
				Value: params.RepoRelDir,
				Short: true,
			},
			{
				Title: "Short-Messages",
				Value: params.Outputs,
				Short: true,
			},
		},
		Footer: fmt.Sprintf("PR #%d • %d commits • %d files changed", params.Number, params.Commits, params.ChangeFileCount),
	}

	_, _, err := api.PostMessage(
		params.SlackChannel,
		slack.MsgOptionAttachments(attachment),
	)

	if err != nil {
		return fmt.Errorf("slack message failed: %w", err)
	}

	fmt.Println("✅ Slack notification sent successfully!")
	return nil
}

// getStatusColor returns color based on status
func getStatusColor(status string) string {
	switch strings.ToLower(status) {
	case NEW:
		return "good" // 초록색
	case PLAN, APPLY:
		return "#0000FF" // 파란색
	case ERROR, FAILED:
		return "danger" // 빨간색
	default: // validate
		return "#808080" // 회색
	}
}

// getStatusEmoji returns emoji based on status
func getStatusEmoji(status string) string {
	switch strings.ToLower(status) {
	case PLAN, APPLY:
		return "✅"
	case VALIDATE:
		return "🔄"
	case ERROR, FAILED:
		return "❌"
	default:
		return "ℹ️"
	}
}

// getCommandTitle returns formatted command title
func getCommandTitle(command string) string {
	switch strings.ToLower(command) {
	case "validate":
		return "Validate"
	case "plan":
		return "Plan"
	case "apply":
		return "Apply"
	default:
		return strings.Title(command)
	}
}
