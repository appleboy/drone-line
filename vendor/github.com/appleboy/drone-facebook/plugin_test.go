package main

import (
	"github.com/stretchr/testify/assert"

	"os"
	"testing"
)

func TestMissingDefaultConfig(t *testing.T) {
	var plugin Plugin

	err := plugin.Exec()

	assert.NotNil(t, err)
}

func TestMissingUserConfig(t *testing.T) {
	plugin := Plugin{
		Config: Config{
			PageToken:   "123456789",
			VerifyToken: "123456789",
		},
	}

	err := plugin.Exec()

	assert.NotNil(t, err)
}

func TestDefaultMessageFormat(t *testing.T) {
	plugin := Plugin{
		Repo: Repo{
			Name:  "go-hello",
			Owner: "appleboy",
		},
		Build: Build{
			Number:  101,
			Status:  "success",
			Link:    "https://github.com/appleboy/go-hello",
			Author:  "Bo-Yi Wu",
			Branch:  "master",
			Message: "update travis",
			Commit:  "e7c4f0a63ceeb42a39ac7806f7b51f3f0d204fd2",
		},
	}

	message := plugin.Message(plugin.Repo, plugin.Build)

	assert.Equal(t, []string{"[success] <https://github.com/appleboy/go-hello> (master)『update travis』by Bo-Yi Wu"}, message)
}

func TestErrorTemplate(t *testing.T) {
	plugin := Plugin{
		Config: Config{
			PageToken:   os.Getenv("FB_PAGE_TOKEN"),
			VerifyToken: os.Getenv("FB_VERIFY_TOKEN"),
			Verify:      false,
			To:          []string{"1234567890"},
			Message:     []string{"file://xxxxx/xxxxx"},
		},
	}

	err := plugin.Exec()
	assert.NotNil(t, err)
}

func TestSendMessage(t *testing.T) {
	plugin := Plugin{
		Repo: Repo{
			Name:  "go-hello",
			Owner: "appleboy",
		},
		Build: Build{
			Number:  101,
			Status:  "success",
			Link:    "https://github.com/appleboy/go-hello",
			Author:  "Bo-Yi Wu",
			Branch:  "master",
			Message: "update travis by drone plugin",
			Commit:  "e7c4f0a63ceeb42a39ac7806f7b51f3f0d204fd2",
			Email:   "test@gmail.com",
		},

		Config: Config{
			PageToken:   os.Getenv("FB_PAGE_TOKEN"),
			VerifyToken: os.Getenv("FB_VERIFY_TOKEN"),
			Verify:      false,
			To:          []string{os.Getenv("FB_TO"), os.Getenv("FB_TO") + ":appleboy@gmail.com", "中文ID", "1234567890"},
			Message:     []string{"Test Facebook Bot From Travis or Local from {{ build.author }}", " "},
			Image:       []string{"https://cdn3.iconfinder.com/data/icons/picons-social/57/16-apple-256.png", "tests/1234.png"},
			Audio:       []string{"https://ia802508.us.archive.org/5/items/testmp3testfile/mpthreetest.mp3", "tests/1234.mp3"},
			Video:       []string{"http://www.sample-videos.com/video/mp4/480/big_buck_bunny_480p_5mb.mp4", "tests/1234.mp4"},
			File:        []string{"http://open.qiniudn.com/where-can-you-use-golang.pdf", "tests/1234.pdf"},
		},
	}

	err := plugin.Exec()
	assert.Nil(t, err)

	// disable message
	plugin.Config.Message = []string{}
	err = plugin.Exec()
	assert.Nil(t, err)
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

func TestParseTo(t *testing.T) {
	input := []string{"0", "1:1@gmail.com", "2:2@gmail.com", "3:3@gmail.com", "4", "5"}

	ids := parseTo(input, "1@gmail.com", false)
	assert.Equal(t, []int64{0, 4, 5, 1}, ids)

	ids = parseTo(input, "1@gmail.com", true)
	assert.Equal(t, []int64{1}, ids)

	ids = parseTo(input, "a@gmail.com", false)
	assert.Equal(t, []int64{0, 4, 5}, ids)

	ids = parseTo(input, "a@gmail.com", true)
	assert.Equal(t, []int64{0, 4, 5}, ids)

	// test empty ids
	ids = parseTo([]string{"", " ", "   "}, "a@gmail.com", true)
	assert.Equal(t, 0, len(ids))
}
