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
			Delimiter:     "::",
			Message:       []string{"Test Line Bot From Travis or Local", " "},
			Image:         []string{"https://cdn3.iconfinder.com/data/icons/picons-social/57/16-apple-128.png"},
			Video:         []string{"http://www.sample-videos.com/video/mp4/480/big_buck_bunny_480p_5mb.mp4"},
			Audio:         []string{"http://feeds-tmp.soundcloud.com/stream/270161326-gotimefm-5-sarah-adams-on-test2doc-and-women-who-go.mp3::2920000", "http://feeds-tmp.soundcloud.com/stream/270161326-gotimefm-5-sarah-adams-on-test2doc-and-women-who-go.mp3"},
			Sticker:       []string{"1::1::100", "1::1"},
			Location:      []string{"竹北體育館::新竹縣竹北市::24.834687::120.993368", "1::1"},
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

func TestConvertImage(t *testing.T) {
	var input string
	var result []string

	input = "http://example.com/1.png"
	result = []string{"http://example.com/1.png", "http://example.com/1.png"}

	assert.Equal(t, result, convertImage(input, "::"))

	input = "http://example.com/1.png::http://example.com/2.png"
	result = []string{"http://example.com/1.png", "http://example.com/2.png"}

	assert.Equal(t, result, convertImage(input, "::"))

	input = "http://example.com/1.png@@http://example.com/2.png"
	result = []string{"http://example.com/1.png", "http://example.com/2.png"}

	assert.Equal(t, result, convertImage(input, "@@"))
}

func TestConvertVideo(t *testing.T) {
	var input string
	var result []string

	input = "http://example.com/1.mp4"
	result = []string{"http://example.com/1.mp4", defaultPreviewImageURL}

	assert.Equal(t, result, convertVideo(input, "::"))

	input = "http://example.com/1.mp4::http://example.com/2.png"
	result = []string{"http://example.com/1.mp4", "http://example.com/2.png"}

	assert.Equal(t, result, convertVideo(input, "::"))

	input = "http://example.com/1.mp4@@http://example.com/2.png"
	result = []string{"http://example.com/1.mp4", "http://example.com/2.png"}

	assert.Equal(t, result, convertVideo(input, "@@"))
}

func TestConvertAudio(t *testing.T) {
	var input string
	var result Audio
	var empty bool

	input = "http://example.com/1.mp3"
	result, empty = convertAudio(input, "::")

	assert.Equal(t, true, empty)
	assert.Equal(t, Audio{}, result)

	// strconv.ParseInt: parsing "我": invalid syntax
	input = "http://example.com/1.mp3::我"
	result, empty = convertAudio(input, "::")

	assert.Equal(t, true, empty)
	assert.Equal(t, Audio{}, result)

	input = "http://example.com/1.mp3::1000"
	result, empty = convertAudio(input, "::")

	assert.Equal(t, false, empty)
	assert.Equal(t, Audio{
		URL:      "http://example.com/1.mp3",
		Duration: 1000,
	}, result)
}

func TestConvertSticker(t *testing.T) {
	var input string
	var result []int
	var empty bool

	input = "1,1"
	result, empty = convertSticker(input, "::")

	assert.Equal(t, true, empty)
	assert.Equal(t, []int{}, result)

	// strconv.ParseInt: parsing "我": invalid syntax
	input = "1::我::100"
	result, empty = convertSticker(input, "::")

	assert.Equal(t, true, empty)
	assert.Equal(t, []int{}, result)

	input = "1::1::100"
	result, empty = convertSticker(input, "::")

	assert.Equal(t, false, empty)
	assert.Equal(t, []int{1, 1, 100}, result)
}

func TestConvertLocation(t *testing.T) {
	var input string
	var result Location
	var empty bool

	input = "1::2::3"
	result, empty = convertLocation(input, "::")

	assert.Equal(t, true, empty)
	assert.Equal(t, Location{}, result)

	// strconv.ParseInt: parsing "測試": invalid syntax
	input = "竹北體育館::新竹縣竹北市::測試::139.704051"
	result, empty = convertLocation(input, "::")

	assert.Equal(t, true, empty)
	assert.Equal(t, Location{}, result)

	// strconv.ParseInt: parsing "測試": invalid syntax
	input = "竹北體育館::新竹縣竹北市::35.661777::測試"
	result, empty = convertLocation(input, "::")

	assert.Equal(t, true, empty)
	assert.Equal(t, Location{}, result)

	input = "竹北體育館::新竹縣竹北市::35.661777::139.704051"
	result, empty = convertLocation(input, "::")

	assert.Equal(t, false, empty)
	assert.Equal(t, Location{
		Title:     "竹北體育館",
		Address:   "新竹縣竹北市",
		Latitude:  float64(35.661777),
		Longitude: float64(139.704051),
	}, result)
}
