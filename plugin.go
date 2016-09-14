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

	log.Print("%v", p.Config)
	ChannelID, err := strconv.ParseInt(p.Config.ChannelID, 10, 64)
	if err != nil {
		log.Fatal("Wrong ChannelID")
	}

	bot, err := linebot.NewClient(ChannelID, p.Config.ChannelSecret, p.Config.MID)
	if err != nil {
		log.Fatal(err)
	}

	to := strings.Split(p.Config.To, ",")

	_, err = bot.SendText(to, p.Config.Message)

	if err != nil {
		log.Println(err)
	}

	return nil
}

func message(repo Repo, build Build) string {
	return fmt.Sprintf("*%s* <%s|%s/%s#%s> (%s) by %s",
		build.Status,
		build.Link,
		repo.Owner,
		repo.Name,
		build.Commit[:8],
		build.Branch,
		build.Author,
	)
}

func fallback(repo Repo, build Build) string {
	return fmt.Sprintf("%s %s/%s#%s (%s) by %s",
		build.Status,
		repo.Owner,
		repo.Name,
		build.Commit[:8],
		build.Branch,
		build.Author,
	)
}

func color(build Build) string {
	switch build.Status {
	case "success":
		return "good"
	case "failure", "error", "killed":
		return "danger"
	default:
		return "warning"
	}
}

func prepend(prefix, s string) string {
	if !strings.HasPrefix(s, prefix) {
		return prefix + s
	}
	return s
}
