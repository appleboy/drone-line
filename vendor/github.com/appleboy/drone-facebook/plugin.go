package main

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/appleboy/drone-facebook/template"
	"github.com/paked/messenger"
)

type (
	// Repo information.
	Repo struct {
		Owner string
		Name  string
	}

	// Build information.
	Build struct {
		Tag      string
		Event    string
		Number   int
		Commit   string
		Message  string
		Branch   string
		Author   string
		Email    string
		Status   string
		Link     string
		Started  float64
		Finished float64
	}

	// Config for the plugin.
	Config struct {
		PageToken   string
		VerifyToken string
		Verify      bool
		MatchEmail  bool
		To          []string
		Message     []string
		Image       []string
		Audio       []string
		Video       []string
		File        []string
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

func parseTo(to []string, authorEmail string, matchEmail bool) []int64 {
	var emails []int64
	var ids []int64
	attachEmail := true

	for _, value := range trimElement(to) {
		idArray := trimElement(strings.Split(value, ":"))

		// check id
		id, err := strconv.ParseInt(idArray[0], 10, 64)
		if err != nil {
			continue
		}

		// check match author email
		if len(idArray) > 1 {
			if email := idArray[1]; email != authorEmail {
				continue
			}

			emails = append(emails, id)
			attachEmail = false
			continue
		}

		ids = append(ids, id)
	}

	if matchEmail == true && attachEmail == false {
		return emails
	}

	for _, value := range emails {
		ids = append(ids, value)
	}

	return ids
}

// Exec executes the plugin.
func (p Plugin) Exec() error {

	if len(p.Config.PageToken) == 0 || len(p.Config.VerifyToken) == 0 || len(p.Config.To) == 0 {
		log.Println("missing facebook config")

		return errors.New("missing facebook config")
	}

	var message []string
	if len(p.Config.Message) > 0 {
		message = p.Config.Message
	} else {
		message = p.Message(p.Repo, p.Build)
	}

	// Create a new messenger client
	client := messenger.New(messenger.Options{
		Verify:      p.Config.Verify,
		Token:       p.Config.PageToken,
		VerifyToken: p.Config.VerifyToken,
	})

	ids := parseTo(p.Config.To, p.Build.Email, p.Config.MatchEmail)

	// send message.
	for _, user := range ids {
		To := messenger.Recipient{
			ID: user,
		}

		// send text notification
		for _, value := range trimElement(message) {
			txt, err := template.RenderTrim(value, p)
			if err != nil {
				return err
			}

			client.Send(To, txt)
		}

		// send image notification
		for _, value := range trimElement(p.Config.Image) {
			client.Attachment(To, messenger.ImageAttachment, value)
		}

		// send audio notification
		for _, value := range trimElement(p.Config.Audio) {
			client.Attachment(To, messenger.AudioAttachment, value)
		}

		// send video notification
		for _, value := range trimElement(p.Config.Video) {
			client.Attachment(To, messenger.VideoAttachment, value)
		}

		// send file notification
		for _, value := range trimElement(p.Config.File) {
			client.Attachment(To, messenger.FileAttachment, value)
		}
	}

	return nil
}

// Message is plugin default message.
func (p Plugin) Message(repo Repo, build Build) []string {
	return []string{fmt.Sprintf("[%s] <%s> (%s)『%s』by %s",
		build.Status,
		build.Link,
		build.Branch,
		build.Message,
		build.Author,
	)}
}
