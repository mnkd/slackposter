package main

import (
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
	slack.PostMessage("Hello, world!")
}
