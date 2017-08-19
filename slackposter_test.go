package slackposter

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/user"
	"path/filepath"
	"testing"
	"time"
)

// TestingConfig is the configuration for slackposter_test.
//
// Note:
// ~/.config/slackposter/config.json
//
// {
//   "slack_webhooks": [
//     {
//       "channel": "#general",
//       "username": "GitHub Status",
//       "icon_emoji": ":octocat:",
//       "webhook_url": "https://hooks.slack.com/services/XXX/YYY/ZZZZ"
//     },
//     ..snip..
// }
//
type TestingConfig struct {
	Webhooks []Config `json:"slack_webhooks"`
}

// loadConfig setup a TestingConfig.
func loadConfig() (TestingConfig, error) {
	var config TestingConfig

	usr, err := user.Current()
	if err != nil {
		fmt.Println("Could not get current user.", err)
		return TestingConfig{}, err
	}

	path := filepath.Join(usr.HomeDir, "/.config/slackposter/config.json")

	str, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("Could not read config.json. ", err)
		return TestingConfig{}, err
	}

	if err := json.Unmarshal(str, &config); err != nil {
		fmt.Println("Failed to unmarshal config.json:", err)
		return config, err
	}

	return config, nil
}

// prepareSlackPoster prepares a SlackPoster instance.
func prepareSlackPoster(t *testing.T) SlackPoster {
	config, err := loadConfig()
	if err != nil {
		t.Fatal(err)
	}
	if len(config.Webhooks) == 0 {
		t.Fatal("len(config.channels) == 0")
	}
	return NewSlackPoster(config.Webhooks[0])
}

func TestPostMessage(t *testing.T) {
	message := "Hello, world!"
	poster := prepareSlackPoster(t)
	err := poster.PostMessage(message)
	if err != nil {
		t.Error(err)
	}
}

func TestPostMessageDryRun(t *testing.T) {
	message := "Hello world. (dry run)"
	poster := prepareSlackPoster(t)
	poster.DryRun = true
	err := poster.PostMessage(message)
	if err != nil {
		t.Error(err)
	}
}

func TestPostPayload(t *testing.T) {
	poster := prepareSlackPoster(t)

	var payload Payload
	payload.Channel = poster.Channel
	payload.Username = "GitHub Status"
	payload.IconEmoji = poster.IconEmoji
	payload.LinkNames = true
	payload.Mrkdwn = true

	statusField := Field{
		Title: "Status",
		Value: "Good",
		Short: true,
	}

	dateString := time.Now().Format("2006-01-02 15:04")
	dateField := Field{
		Title: "Date",
		Value: dateString,
		Short: true,
	}

	attachment := Attachment{
		Fallback: "GitHub Status: Good - https://status.github.com",
		Text:     "<https://status.github.com/|GitHub Status> : *Good*",
		Color:    "good",
		Fields:   []Field{statusField, dateField},
		MrkdwnIn: []string{"text"},
	}

	payload.Attachments = []Attachment{attachment}

	err := poster.PostPayload(payload)
	if err != nil {
		t.Error(err)
	}
}
