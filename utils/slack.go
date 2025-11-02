package utils

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/slack-go/slack"
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

	placeHolder := fmt.Sprintf("⚠️ Drift detected in %d project(s)\n", countNonEmptyOutputs(tfOutputs))
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
