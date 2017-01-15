package main

import (
	"time"

	"github.com/mnkd/slackposter"
)

func main() {
	config := slackposter.Config{
		"#your-channel",
		":octocat:",
		"GitHub | Status",
		"https://hooks.slack.com/services/xxx/xxx/xxx",
	}

	slack := slackposter.NewSlack(config)

	payload := slack.NewPayload()
	payload.Mrkdwn = true

	statusField := slackposter.Field{
		Title: "Status",
		Value: "Good",
		Short: true,
	}

	dateString := time.Now().Format("2006-01-02 15:04")
	dateField := slackposter.Field{
		Title: "Date",
		Value: dateString,
		Short: true,
	}

	attachment := slackposter.Attachment{
		Fallback: "GitHub Status: Good - https://status.github.com",
		Text:     "<https://status.github.com/|GitHub Status> : *Good*",
		Color:    "good",
		Fields:   []slackposter.Field{statusField, dateField},
		MrkdwnIn: []string{"text"},
	}

	payload.Attachments = []slackposter.Attachment{attachment}

	slack.PostPayload(payload)
}
