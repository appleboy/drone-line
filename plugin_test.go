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
	var plugin Plugin

	plugin.Config.ChannelID = "123456789"
	plugin.Config.ChannelSecret = "test wrong id"
	plugin.Config.MID = "test wrong id"

	err := plugin.Exec()

	assert.NotNil(t, err)
}

func TestSendTextError(t *testing.T) {
	var plugin Plugin

	plugin.Repo.Name = "octocat"
	plugin.Repo.Owner = "hello-world"
	plugin.Build.Number = 100
	plugin.Build.Status = "success"
	plugin.Build.Link = "http://github.com/appleboy/go-hello"
	plugin.Build.Author = "appleboy"
	plugin.Build.Branch = "master"
	plugin.Build.Commit = "7fd1a60b01f91b314f59955a4e4d4e80d8edf11d"
	plugin.Config.ChannelID = "1465486347"
	plugin.Config.ChannelSecret = "ChannelSecret"
	plugin.Config.MID = "MID"
	plugin.Config.To = []string{"1234567890"}
	plugin.Config.Message = "Test"

	// enable message
	err := plugin.Exec()
	assert.NotNil(t, err)

	// disable message
	plugin.Config.Message = ""
	err = plugin.Exec()
	assert.NotNil(t, err)
}

func TestDefaultMessage(t *testing.T) {
	var plugin Plugin

	plugin.Repo.Name = "go-hello"
	plugin.Repo.Owner = "appleboy"
	plugin.Build.Number = 100
	plugin.Build.Status = "success"
	plugin.Build.Link = "http://github.com/appleboy/go-hello"
	plugin.Build.Author = "appleboy"
	plugin.Build.Branch = "master"
	plugin.Build.Commit = "7fd1a60b01f91b314f59955a4e4d4e80d8edf11d"
	plugin.Config.ChannelID = "1465486347"
	plugin.Config.ChannelSecret = "ChannelSecret"
	plugin.Config.MID = "MID"
	plugin.Config.To = []string{"1234567890"}
	plugin.Config.Message = "Test"

	message := plugin.Message(plugin.Repo, plugin.Build)

	assert.Equal(t, "[success] <http://github.com/appleboy/go-hello|appleboy/go-hello#7fd1a60b> (master) by appleboy", message)
}
