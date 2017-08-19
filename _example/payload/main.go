package main

import (
	"time"

	slack "github.com/mnkd/slackposter"
)

func main() {
	config := slack.Config{
		Channel:    "#your-channel",
		IconEmoji:  ":octocat:",
		Username:   "GitHub | Status",
		WebhookUrl: "https://hooks.slack.com/services/xxx/xxx/xxx",
	}

	poster := slack.NewSlackPoster(config)

	payload := poster.NewPayload()
	payload.Mrkdwn = true

	statusField := slack.Field{
		Title: "Status",
		Value: "Good",
		Short: true,
	}

	dateString := time.Now().Format("2006-01-02 15:04")
	dateField := slack.Field{
		Title: "Date",
		Value: dateString,
		Short: true,
	}

	attachment := slack.Attachment{
		Fallback: "GitHub Status: Good - https://status.github.com",
		Text:     "<https://status.github.com/|GitHub Status> : *Good*",
		Color:    "good",
		Fields:   []slack.Field{statusField, dateField},
		MrkdwnIn: []string{"text"},
	}

	payload.Attachments = []slack.Attachment{attachment}

	poster.PostPayload(payload)
}
