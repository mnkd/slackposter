# skackposter
[![Godoc Reference](https://godoc.org/github.com/mnkd/slackposter?status.svg)](https://godoc.org/github.com/mnkd/slackposter)

Post a payload to your Slack incoming webhook.

## Usage

You can send a simple message.

```go
package main

import (
	"github.com/mnkd/slackposter"
)

func main() {
    config := slackposter.Config{
        "#channel_name or channel_id",
        ":ghost:",
        "Ghost",
        "https://hooks.slack.com/services/xxxx/xxx/xxx",
    }

    slack := slackposter.NewSlack(config)
    slack.PostMessage("Hello world!")
}
```

![message](_examples/message/message.png)

You can send a customized message (payload).

```go
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
		WebhookURL: "https://hooks.slack.com/services/xxx/xxx/xxx",
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
```

![payload](_examples/payload/payload.png)

I would recommend you read
[Attaching content and links to messages | Slack](https://api.slack.com/docs/message-attachments).

## Example Apps
### [github-status](https://github.com/mnkd/github-status)
* Notify GitHub Site Status to Slack incoming webhook.

![github-status](https://github.com/mnkd/github-status/raw/master/images/slack.png)

### [prnotify](https://github.com/mnkd/prnotify)
* Notify GitHub pull requests to Slack incoming webhook.

### [qiiotd](https://github.com/mnkd/qiiotd)
* Qiita:Team (Qiita) の n 年前の今日の記事を Slack の Incoming webhook へ post する。

![qiiotd](https://github.com/mnkd/qiiotd/raw/master/images/slack.png)
