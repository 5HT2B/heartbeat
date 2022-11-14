package main

import (
	"encoding/json"
	"github.com/valyala/fasthttp"
	"log"
	"net/url"
)

var (
	liveUrl      = ""                 // Set in main.go after env is parsed
	webhookUrl   = ""                 // Set in main.go after env is parsed
	webhookLevel = WebhookLevelSimple // Set in main.go after env is parsed
)

type EmbedColor int

const (
	EmbedColorGreen  EmbedColor = 4388248
	EmbedColorOrange EmbedColor = 14587196
	EmbedColorBlue   EmbedColor = 6591981
)

type WebhookLevel int64

const (
	WebhookLevelAll         WebhookLevel = iota // Log each beat + the below
	WebhookLevelSimple                          // Log new devices + the below
	WebhookLevelLongAbsence                     // Log absences longer than 1h
	WebhookLevelNone                            // Don't log anything
)

func (wl WebhookLevel) String() string {
	switch wl {
	case WebhookLevelAll:
		return "ALL"
	case WebhookLevelSimple:
		return "SIMPLE"
	case WebhookLevelLongAbsence:
		return "LONG_ABSENCE"
	case WebhookLevelNone:
		return "NONE"
	default:
		return WebhookLevelSimple.String()
	}
}

func PostMessage(title, description string, color EmbedColor, level WebhookLevel) {
	if len(webhookUrl) == 0 || len(liveUrl) == 0 {
		return // hasn't been set in config/.env
	}

	// Log level too low
	if webhookLevel > level {
		log.Printf("Returning %s / %s\n", webhookLevel, level)
		return
	}

	type discordAuthor struct {
		Name    string `json:"name,omitempty"`
		Url     string `json:"url,omitempty"`
		IconUrl string `json:"icon_url,omitempty"`
	}
	type discordThumbnail struct {
		Url string `json:"url,omitempty"`
	}
	type discordField struct {
		Title       string `json:"name,omitempty"`
		Description string `json:"value,omitempty"`
	}
	type discordEmbed struct {
		Title       string           `json:"title,omitempty"`
		Description string           `json:"description,omitempty"`
		Color       EmbedColor       `json:"color,omitempty"`
		Author      discordAuthor    `json:"author,omitempty"`
		Thumbnail   discordThumbnail `json:"thumbnail,omitempty"`
		Fields      []discordField   `json:"fields,omitempty"`
	}
	type discordMessage struct {
		Content   string         `json:"content,omitempty"`
		AvatarUrl string         `json:"avatar_url,omitempty"`
		Username  string         `json:"username,omitempty"`
		Embeds    []discordEmbed `json:"embeds,omitempty"`
	}

	host := serverName
	if u, err := url.Parse(liveUrl); err == nil {
		host = u.Host
	}
	avatar := liveUrl + "/favicon.png"
	author := discordAuthor{host, liveUrl, avatar}
	message := discordMessage{AvatarUrl: avatar, Username: host, Embeds: []discordEmbed{{
		Title:       title,
		Description: description,
		Color:       color,
		Author:      author,
	}}}

	j, err := json.Marshal(message)

	if err != nil {
		if *debug {
			log.Printf("Error creating webhook message json: %s", err)
		}
		return // If this method fails, it doesn't matter to users
	}

	if *debug {
		log.Printf("Posting webhook json: %s", j)
	}

	req := fasthttp.AcquireRequest()
	req.SetBody(j)
	req.Header.SetMethod(fasthttp.MethodPost)
	req.Header.SetContentType(jsonMime)
	req.SetRequestURI(webhookUrl)
	res := fasthttp.AcquireResponse()
	if err := fasthttp.Do(req, res); err != nil {
		fasthttp.ReleaseRequest(req)
		log.Printf("Error posting webhook: %s", err)
		return
	}
	fasthttp.ReleaseRequest(req)
	resBody := res.Body()
	if *debug && len(resBody) > 0 {
		log.Printf("Webhook response: %s", resBody)
	}
	fasthttp.ReleaseResponse(res) // When done with resBody
}
