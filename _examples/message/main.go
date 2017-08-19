package main

import (
	slack "github.com/mnkd/slackposter"
)

func main() {
	config := slack.Config{
		Channel:    "#your-channel",
		IconEmoji:  ":octocat:",
		Username:   "Octocat",
		WebhookURL: "https://hooks.slack.com/services/xxx/xxx/xxx",
	}
	poster := slack.NewSlackPoster(config)
	poster.PostMessage("Hello, world!")
}
