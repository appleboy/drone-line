package main

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/line/line-bot-sdk-go/linebot"
)

type (
	// Repo information.
	Repo struct {
		Owner string
		Name  string
	}

	// Build information.
	Build struct {
		Event  string
		Number int
		Commit string
		Branch string
		Author string
		Status string
		Link   string
	}

	// Config for the plugin.
	Config struct {
		ChannelID     string
		ChannelSecret string
		MID           string
		To            []string
		Message       []string
		Image         []string
		Video         []string
	}

	// Plugin values.
	Plugin struct {
		Repo   Repo
		Build  Build
		Config Config
	}
)

func trimElement(keys []string) []string {
	var newKeys []string

	for _, value := range keys {
		value = strings.Trim(value, " ")
		if len(value) == 0 {
			continue
		}
		newKeys = append(newKeys, value)
	}

	return newKeys
}

// Exec executes the plugin.
func (p Plugin) Exec() error {

	if len(p.Config.ChannelID) == 0 || len(p.Config.ChannelSecret) == 0 || len(p.Config.MID) == 0 {
		log.Println("missing line bot config")

		return errors.New("missing line bot config")
	}

	ChannelID, err := strconv.ParseInt(p.Config.ChannelID, 10, 64)
	if err != nil {
		log.Println("wrong channel id")

		return err
	}

	bot, _ := linebot.NewClient(ChannelID, p.Config.ChannelSecret, p.Config.MID)

	if len(p.Config.To) == 0 {
		log.Println("missing line user config")

		return errors.New("missing line user config")
	}

	var message []string
	if len(p.Config.Message) > 0 {
		message = p.Config.Message
	} else {
		message = p.Message(p.Repo, p.Build)
	}

	// New multiple request instance
	line := bot.NewMultipleMessage()

	// check message array.
	for _, value := range trimElement(message) {
		line.AddText(value)
	}

	// check image array.
	for _, value := range trimElement(p.Config.Image) {
		line.AddImage(value, value)
	}

	// check video array.
	for _, value := range trimElement(p.Config.Video) {
		line.AddVideo(value, "https://cdn4.iconfinder.com/data/icons/miu/24/device-camera-recorder-video-glyph-256.png")
	}

	_, err = line.Send(p.Config.To)

	if err != nil {
		log.Println(err.Error())

		return err
	}

	return nil
}

// Message is line default message.
func (p Plugin) Message(repo Repo, build Build) []string {
	return []string{fmt.Sprintf("[%s] <%s> (%s) by %s",
		build.Status,
		build.Link,
		build.Branch,
		build.Author,
	)}
}
