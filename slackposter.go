package slackposter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// SlackPoster represents a poster for Slack incoming webhook.
// SlackPoster must be created by calling NewSlackPoster().
type SlackPoster struct {
	Channel    string
	DryRun     bool
	IconEmoji  string
	Username   string
	WebhookURL string
}

// Payload is represented a Slack messsage format.
//
// Refer to:
// Attaching content and links to messages | Slack
// https://api.slack.com/docs/message-attachments
//
type Payload struct {
	Attachments []Attachment `json:"attachments"`
	Channel     string       `json:"channel"`
	IconEmoji   string       `json:"icon_emoji"`
	LinkNames   bool         `json:"link_names"`
	Mrkdwn      bool         `json:"mrkdwn"`
	Text        string       `json:"text"`
	Username    string       `json:"username"`
}

// Field represents a `field`.
type Field struct {
	Short bool   `json:"short"`
	Title string `json:"title"`
	Value string `json:"value"`
}

// Attachment represents a `attachment`.
type Attachment struct {
	AuthorIcon string   `json:"author_icon"`
	AuthorLink string   `json:"author_link"`
	AuthorName string   `json:"author_name"`
	Color      string   `json:"color"`
	Fallback   string   `json:"fallback"`
	Fields     []Field  `json:"fields"`
	FooterIcon string   `json:"footer_icon"`
	Footer     string   `json:"footer"`
	ImageURL   string   `json:"image_url"`
	MrkdwnIn   []string `json:"mrkdwn_in"`
	Pretext    string   `json:"pretext"`
	Text       string   `json:"text"`
	ThumbURL   string   `json:"thumb_url"`
	TitleLink  string   `json:"title_link"`
	Title      string   `json:"title"`
	Ts         int64    `json:"ts"`
}

// Config is the configuration of SlackPoster.
type Config struct {
	Channel    string `json:"channel"`
	IconEmoji  string `json:"icon_emoji"`
	Username   string `json:"username"`
	WebhookURL string `json:"webhook_url"`
}

// NewSlackPoster create a new SlackPoster given a Config.
func NewSlackPoster(c Config) SlackPoster {
	return SlackPoster{
		Channel:    c.Channel,
		IconEmoji:  c.IconEmoji,
		Username:   c.Username,
		WebhookURL: c.WebhookURL,
	}
}

// NewPayload create a new Payload.
func (sp SlackPoster) NewPayload() Payload {
	return Payload{
		Channel:   sp.Channel,
		Username:  sp.Username,
		IconEmoji: sp.IconEmoji,
		LinkNames: true,
	}
}

// AppendField append a Field to Attachments.
func (payload *Payload) AppendField(field Field, attachmentIndex int) {
	if attachmentIndex >= len(payload.Attachments) {
		return
	}
	payload.Attachments[attachmentIndex].Fields = append(payload.Attachments[attachmentIndex].Fields, field)
}

// PostPayload posts a payload to WebhookURL.
func (sp SlackPoster) PostPayload(payload Payload) error {
	return sp.post(payload)
}

// PostMessage posts a message to WebhookURL.
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
	request, err := http.NewRequest("POST", sp.WebhookURL, body)
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
