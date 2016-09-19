package main

import (
	"github.com/stretchr/testify/assert"

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
			Message:       "Test",
		},
	}

	// enable message
	err := plugin.Exec()
	assert.NotNil(t, err)

	// disable message
	plugin.Config.Message = ""
	err = plugin.Exec()
	assert.NotNil(t, err)
}

func TestDefaultMessage(t *testing.T) {
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
			Message:       "Test",
		},
	}

	message := plugin.Message(plugin.Repo, plugin.Build)

	assert.Equal(t, "[success] <https://github.com/appleboy/go-hello|appleboy/go-hello#e7c4f0a6> (master) by Bo-Yi Wu", message)
}
