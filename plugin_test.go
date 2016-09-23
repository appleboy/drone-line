package main

import (
	"github.com/stretchr/testify/assert"

	"os"
	"testing"
)

func TestMissingLineConfig(t *testing.T) {
	var plugin Plugin

	err := plugin.Exec()

	assert.NotNil(t, err)
}

func TestWrongChannelID(t *testing.T) {
	var plugin Plugin

	plugin.Config.ChannelID = "test wrong id"
	plugin.Config.ChannelSecret = "test wrong id"
	plugin.Config.MID = "test wrong id"

	err := plugin.Exec()

	assert.NotNil(t, err)
}

func TestMissingUserConfig(t *testing.T) {
	plugin := Plugin{
		Config: Config{
			ChannelID:     "123456789",
			ChannelSecret: "test wrong id",
			MID:           "test wrong id",
		},
	}

	err := plugin.Exec()

	assert.NotNil(t, err)
}

func TestSendTextError(t *testing.T) {
	plugin := Plugin{
		Repo: Repo{
			Name:  "go-hello",
			Owner: "appleboy",
		},
		Build: Build{
			Number: 101,
			Status: "success",
			Link:   "https://github.com/appleboy/go-hello",
			Author: "Bo-Yi Wu",
			Branch: "master",
			Commit: "e7c4f0a63ceeb42a39ac7806f7b51f3f0d204fd2",
		},
		Config: Config{
			ChannelID:     "1465486347",
			ChannelSecret: "ChannelSecret",
			MID:           "MID",
			To:            []string{"1234567890"},
			Message:       []string{"Test"},
		},
	}

	// enable message
	err := plugin.Exec()
	assert.NotNil(t, err)

	// disable message
	plugin.Config.Message = []string{}
	err = plugin.Exec()
	assert.NotNil(t, err)
}

func TestDefaultMessageFormat(t *testing.T) {
	plugin := Plugin{
		Repo: Repo{
			Name:  "go-hello",
			Owner: "appleboy",
		},
		Build: Build{
			Number: 101,
			Status: "success",
			Link:   "https://github.com/appleboy/go-hello",
			Author: "Bo-Yi Wu",
			Branch: "master",
			Commit: "e7c4f0a63ceeb42a39ac7806f7b51f3f0d204fd2",
		},
	}

	message := plugin.Message(plugin.Repo, plugin.Build)

	assert.Equal(t, []string{"[success] <https://github.com/appleboy/go-hello> (master) by Bo-Yi Wu"}, message)
}

func TestErrorSendMessage(t *testing.T) {
	plugin := Plugin{
		Repo: Repo{
			Name:  "go-hello",
			Owner: "appleboy",
		},
		Build: Build{
			Number: 101,
			Status: "success",
			Link:   "https://github.com/appleboy/go-hello",
			Author: "Bo-Yi Wu",
			Branch: "master",
			Commit: "e7c4f0a63ceeb42a39ac7806f7b51f3f0d204fd2",
		},

		Config: Config{
			ChannelID:     os.Getenv("LINE_CHANNEL_ID"),
			ChannelSecret: os.Getenv("LINE_CHANNEL_SECRET"),
			MID:           os.Getenv("LINE_MID"),
			To:            []string{os.Getenv("LINE_TO")},
			Message:       []string{"Test Line Bot From Travis or Local", " "},
			Image:         []string{"https://cdn3.iconfinder.com/data/icons/picons-social/57/16-apple-128.png"},
			Video:         []string{"http://www.sample-videos.com/video/mp4/480/big_buck_bunny_480p_5mb.mp4"},
		},
	}

	err := plugin.Exec()
	// error message: Your ip address [xxx.xxx.xxx.xxx] is not allowed to access this API.
	// Please add your IP to the IP whitelist in the developer center.
	assert.NotNil(t, err)
}

func TestTrimElement(t *testing.T) {
	var input, result []string

	input = []string{"1", "     ", "3"}
	result = []string{"1", "3"}

	assert.Equal(t, result, trimElement(input))

	input = []string{"1", "2"}
	result = []string{"1", "2"}

	assert.Equal(t, result, trimElement(input))
}
