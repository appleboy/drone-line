package main

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/appleboy/drone-facebook/template"
	"github.com/line/line-bot-sdk-go/linebot"
)

const defaultPreviewImageURL = "https://cdn4.iconfinder.com/data/icons/miu/24/device-camera-recorder-video-glyph-256.png"

type (
	// Repo information.
	Repo struct {
		Owner string
		Name  string
	}

	// Build information.
	Build struct {
		Event    string
		Number   int
		Commit   string
		Branch   string
		Author   string
		Message  string
		Status   string
		Link     string
		Started  float64
		Finished float64
	}

	// Config for the plugin.
	Config struct {
		ChannelToken  string
		ChannelSecret string
		To            []string
		Delimiter     string
		Message       []string
		Image         []string
		Video         []string
		Audio         []string
		Sticker       []string
		Location      []string
	}

	// Plugin values.
	Plugin struct {
		Repo   Repo
		Build  Build
		Config Config
	}

	// Audio format
	Audio struct {
		URL      string
		Duration int
	}

	// Location format
	Location struct {
		Title     string
		Address   string
		Latitude  float64
		Longitude float64
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

func convertImage(value, delimiter string) []string {
	values := trimElement(strings.Split(value, delimiter))

	if len(values) < 2 {
		values = append(values, values[0])
	}

	return values
}

func convertVideo(value, delimiter string) []string {
	values := trimElement(strings.Split(value, delimiter))

	if len(values) < 2 {
		values = append(values, defaultPreviewImageURL)
	}

	return values
}

func convertAudio(value, delimiter string) (Audio, bool) {
	values := trimElement(strings.Split(value, delimiter))

	if len(values) < 2 {
		return Audio{}, true
	}

	duration, err := strconv.Atoi(values[1])

	if err != nil {
		log.Println(err.Error())
		return Audio{}, true
	}

	return Audio{
		URL:      values[0],
		Duration: duration,
	}, false
}

func convertSticker(value, delimiter string) ([]string, bool) {
	values := trimElement(strings.Split(value, delimiter))

	if len(values) < 2 {
		return []string{}, true
	}

	return values, false
}

// func convertLocation(value, delimiter string) (Location, bool) {
// 	var latitude, longitude float64
// 	var err error
// 	values := trimElement(strings.Split(value, delimiter))

// 	if len(values) < 4 {
// 		return Location{}, true
// 	}

// 	latitude, err = strconv.ParseFloat(values[2], 64)

// 	if err != nil {
// 		log.Println(err.Error())
// 		return Location{}, true
// 	}

// 	longitude, err = strconv.ParseFloat(values[3], 64)

// 	if err != nil {
// 		log.Println(err.Error())
// 		return Location{}, true
// 	}

// 	return Location{
// 		Title:     values[0],
// 		Address:   values[1],
// 		Latitude:  latitude,
// 		Longitude: longitude,
// 	}, false
// }

// Exec executes the plugin.
func (p Plugin) Exec() error {

	if len(p.Config.ChannelToken) == 0 || len(p.Config.ChannelSecret) == 0 {
		log.Println("missing line bot config")

		return errors.New("missing line bot config")
	}

	bot, _ := linebot.New(p.Config.ChannelSecret, p.Config.ChannelToken)

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

	// Initial messages array.
	var messages []linebot.Message

	for _, value := range trimElement(message) {
		txt, err := template.RenderTrim(value, p)
		if err != nil {
			return err
		}

		messages = append(messages, linebot.NewTextMessage(txt))
	}

	// Add image message
	for _, value := range trimElement(p.Config.Image) {
		values := convertImage(value, p.Config.Delimiter)

		messages = append(messages, linebot.NewImageMessage(values[0], values[1]))
	}

	// Add image message.
	for _, value := range trimElement(p.Config.Video) {
		values := convertVideo(value, p.Config.Delimiter)

		messages = append(messages, linebot.NewVideoMessage(values[0], values[1]))
	}

	// Add Audio message.
	for _, value := range trimElement(p.Config.Audio) {
		audio, empty := convertAudio(value, p.Config.Delimiter)

		if empty == true {
			continue
		}

		messages = append(messages, linebot.NewAudioMessage(audio.URL, audio.Duration))
	}

	// Add Sticker message.
	for _, value := range trimElement(p.Config.Sticker) {
		sticker, empty := convertSticker(value, p.Config.Delimiter)

		if empty == true {
			continue
		}

		messages = append(messages, linebot.NewStickerMessage(sticker[0], sticker[1]))
	}

	// // check Location array.
	// for _, value := range trimElement(p.Config.Location) {
	// 	location, empty := convertLocation(value, p.Config.Delimiter)

	// 	if empty == true {
	// 		continue
	// 	}

	// 	line.AddLocation(location.Title, location.Address, location.Latitude, location.Longitude)
	// }

	// send message to user
	for _, id := range trimElement(p.Config.To) {
		if _, err := bot.PushMessage(id, messages...).Do(); err != nil {
			log.Println(err.Error())
		}
	}

	return nil
}

// Message is line default message.
func (p Plugin) Message(repo Repo, build Build) []string {
	return []string{fmt.Sprintf("[%s] <%s> (%s)『%s』by %s",
		build.Status,
		build.Link,
		build.Branch,
		build.Message,
		build.Author,
	)}
}
