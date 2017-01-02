package slackposter

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/user"
	"path/filepath"
	"testing"
	"time"

	slack "github.com/m-nakada/slackposter"
)

type Config struct {
	Channels []slack.Config `json:"channels"`
}

// Note:
// .config/slackposter/config.json
//
// {
//   "channels": [
//     {
//       "channel": "#general",
//       "username": "GitHub Status",
//       "icon_emoji": ":octocat:",
//       "webhook_url": "https://hooks.slack.com/services/XXX/YYY/ZZZZ"
//     },
//     ..snip..
// }

func NewConfig() (Config, error) {
	var config Config

	usr, err := user.Current()
	if err != nil {
		fmt.Println("Could not get current user.", err)
		return config, err
	}
	path := filepath.Join(usr.HomeDir, "/.config/slackposter/config.json")

	str, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("Could not read config.json. ", err)
		return config, err
	}

	if err := json.Unmarshal(str, &config); err != nil {
		fmt.Println("JSON Unmarshal Error:", err)
		return config, err
	}

	return config, nil
}

func NewSlack(t *testing.T) slack.Slack {
	config, err := NewConfig()
	if err != nil {
		t.Fatal(err)
	}
	if len(config.Channels) == 0 {
		t.Fatal("len(config.channels) == 0")
	}
	slackposter := slack.NewSlack(config.Channels[0])
	return slackposter
}

func TestPostMessage(t *testing.T) {
	message := "Hello world."
	slk := NewSlack(t)
	err := slk.PostMessage(message)
	if err != nil {
		t.Error(err)
	}
}

func TestPostMessageDryRun(t *testing.T) {
	message := "Hello world. (dry run)"
	slk := NewSlack(t)
	slk.DryRun = true
	err := slk.PostMessage(message)
	if err != nil {
		t.Error(err)
	}
}

func TestPostPayload(t *testing.T) {
	slk := NewSlack(t)

	var payload slack.Payload
	payload.Channel = slk.Channel
	payload.Username = "GitHub Status"
	payload.IconEmoji = slk.IconEmoji
	payload.LinkNames = true
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

	err := slk.PostPayload(payload)
	if err != nil {
		t.Error(err)
	}
}
