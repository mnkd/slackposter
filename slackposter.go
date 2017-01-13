package slackposter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type Slack struct {
	Channel    string
	DryRun     bool
	IconEmoji  string
	Username   string
	WebhookUrl string
}

// Refer to:
// Attaching content and links to messages | Slack
// https://api.slack.com/docs/message-attachments

type Payload struct {
	Attachments []Attachment `json:"attachments"`
	Channel     string       `json:"channel"`
	IconEmoji   string       `json:"icon_emoji"`
	LinkNames   bool         `json:"link_names"`
	Mrkdwn      bool         `json:"mrkdwn"`
	Text        string       `json:"text"`
	Username    string       `json:"username"`
}

type Field struct {
	Short bool   `json:"short"`
	Title string `json:"title"`
	Value string `json:"value"`
}

type Attachment struct {
	AuthorIcon string   `json:"author_icon"`
	AuthorLink string   `json:"author_link"`
	AuthorName string   `json:"author_name"`
	Color      string   `json:"color"`
	Fallback   string   `json:"fallback"`
	Fields     []Field  `json:"fields"`
	MrkdwnIn   []string `json:"mrkdwn_in"`
	Text       string   `json:"text"`
	ThumbUrl   string   `json:"thumb_url"`
}

type Config struct {
	Channel    string `json:"channel"`
	IconEmoji  string `json:"icon_emoji"`
	Username   string `json:"username"`
	WebhookUrl string `json:"webhook_url"`
}

func NewSlack(config Config) Slack {
	slack := Slack{
		Channel:    config.Channel,
		IconEmoji:  config.IconEmoji,
		Username:   config.Username,
		WebhookUrl: config.WebhookUrl,
	}
	return slack
}

func (slack Slack) NewPayload() Payload {
	return Payload{
		Channel:   slack.Channel,
		Username:  slack.Username,
		IconEmoji: slack.IconEmoji,
		LinkNames: true,
	}
}

func (payload *Payload) AppendField(field Field, attachmentIndex int) {
	if attachmentIndex >= len(payload.Attachments) {
		return
	}
	payload.Attachments[attachmentIndex].Fields = append(payload.Attachments[attachmentIndex].Fields, field)
}

func (slack Slack) PostPayload(payload Payload) error {
	return slack.post(payload)
}

func (slack Slack) PostMessage(message string) error {
	var payload Payload
	payload.Channel = slack.Channel
	payload.Username = slack.Username
	payload.IconEmoji = slack.IconEmoji
	payload.LinkNames = true
	payload.Text = message
	return slack.post(payload)
}

func (slack Slack) post(payload Payload) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	if slack.DryRun {
		fmt.Fprintln(os.Stdout, "**dry run**\njson:\n", string(data))
		return nil
	}

	// Create request
	body := bytes.NewBuffer(data)
	request, err := http.NewRequest("POST", slack.WebhookUrl, body)
	request.Header.Add("Content-Type", "application/json; charset=utf-8")
	if err != nil {
		return err
	}

	// Invoke Request
	client := &http.Client{}
	_, err = client.Do(request)
	if err != nil {
		return err
	}

	return nil
}
