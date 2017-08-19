package slackposter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type SlackPoster struct {
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
	FooterIcon string   `json:"footer_icon"`
	Footer     string   `json:"footer"`
	ImageUrl   string   `json:"image_url"`
	MrkdwnIn   []string `json:"mrkdwn_in"`
	Pretext    string   `json:"pretext"`
	Text       string   `json:"text"`
	ThumbUrl   string   `json:"thumb_url"`
	TitleLink  string   `json:"title_link"`
	Title      string   `json:"title"`
	Ts         int64    `json:"ts"`
}

type Config struct {
	Channel    string `json:"channel"`
	IconEmoji  string `json:"icon_emoji"`
	Username   string `json:"username"`
	WebhookUrl string `json:"webhook_url"`
}

func NewSlackPoster(c Config) SlackPoster {
	return SlackPoster{
		Channel:    c.Channel,
		IconEmoji:  c.IconEmoji,
		Username:   c.Username,
		WebhookUrl: c.WebhookUrl,
	}
}

func (sp SlackPoster) NewPayload() Payload {
	return Payload{
		Channel:   sp.Channel,
		Username:  sp.Username,
		IconEmoji: sp.IconEmoji,
		LinkNames: true,
	}
}

func (payload *Payload) AppendField(field Field, attachmentIndex int) {
	if attachmentIndex >= len(payload.Attachments) {
		return
	}
	payload.Attachments[attachmentIndex].Fields = append(payload.Attachments[attachmentIndex].Fields, field)
}

func (sp SlackPoster) PostPayload(payload Payload) error {
	return sp.post(payload)
}

func (sp SlackPoster) PostMessage(message string) error {
	var payload Payload
	payload.Channel = sp.Channel
	payload.Username = sp.Username
	payload.IconEmoji = sp.IconEmoji
	payload.LinkNames = true
	payload.Text = message
	return sp.post(payload)
}

func (sp SlackPoster) post(payload Payload) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	if sp.DryRun {
		fmt.Fprintln(os.Stdout, "**dry run**\njson:\n", string(data))
		return nil
	}

	// Create request
	body := bytes.NewBuffer(data)
	request, err := http.NewRequest("POST", sp.WebhookUrl, body)
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
