package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/line/line-bot-sdk-go/linebot"
)

type (
	Repo struct {
		Owner string
		Name  string
	}

	Build struct {
		Event  string
		Number int
		Commit string
		Branch string
		Author string
		Status string
		Link   string
	}

	Config struct {
		ChannelID     string
		ChannelSecret string
		MID           string
		To            string
		Message       string
	}

	Plugin struct {
		Repo   Repo
		Build  Build
		Config Config
	}
)

func (p Plugin) Exec() error {

	ChannelID, err := strconv.ParseInt(p.Config.ChannelID, 10, 64)
	if err != nil {
		log.Fatal("Wrong ChannelID")
	}

	bot, err := linebot.NewClient(ChannelID, p.Config.ChannelSecret, p.Config.MID)
	if err != nil {
		log.Fatal(err)
	}

	to := strings.Split(p.Config.To, ",")

	var message string
	if p.Config.Message != "" {
		message = p.Config.Message
	} else {
		message = p.Message(p.Repo, p.Build)
	}

	_, err = bot.SendText(to, message)

	if err != nil {
		log.Println(err)
	}

	return nil
}

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
