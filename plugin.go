package main

import (
	"errors"
	"fmt"
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
		To            string
		Message       string
	}

	// Plugin values.
	Plugin struct {
		Repo   Repo
		Build  Build
		Config Config
	}
)

// Exec executes the plugin.
func (p Plugin) Exec() error {

	if len(p.Config.ChannelID) == 0 || len(p.Config.ChannelSecret) == 0 || len(p.Config.MID) == 0 {
		return errors.New("missing line bot config")
	}

	ChannelID, err := strconv.ParseInt(p.Config.ChannelID, 10, 64)
	if err != nil {
		return err
	}

	bot, _ := linebot.NewClient(ChannelID, p.Config.ChannelSecret, p.Config.MID)

	to := strings.Split(p.Config.To, ",")

	var message string
	if p.Config.Message != "" {
		message = p.Config.Message
	} else {
		message = p.Message(p.Repo, p.Build)
	}

	_, err = bot.SendText(to, message)

	if err != nil {
		return err
	}

	return nil
}

// Message is line default message.
func (p Plugin) Message(repo Repo, build Build) string {
	return fmt.Sprintf("[%s] <%s|%s/%s#%s> (%s) by %s",
		build.Status,
		build.Link,
		repo.Owner,
		repo.Name,
		build.Commit[:8],
		build.Branch,
		build.Author,
	)
}
