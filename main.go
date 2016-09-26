package main

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/urfave/cli"
)

// Version for command line
var Version string

func main() {
	app := cli.NewApp()
	app.Name = "line plugin"
	app.Usage = "line plugin"
	app.Action = run
	app.Version = Version
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "channel.id",
			Usage:  "line channel id",
			EnvVar: "PLUGIN_CHANNEL_ID,LINE_CHANNEL_ID",
		},
		cli.StringFlag{
			Name:   "channel.secret",
			Usage:  "line channel secret",
			EnvVar: "PLUGIN_CHANNEL_SECRET,LINE_CHANNEL_SECRET",
		},
		cli.StringFlag{
			Name:   "mid",
			Usage:  "line mid",
			EnvVar: "PLUGIN_MID,LINE_MID",
		},
		cli.StringSliceFlag{
			Name:   "to",
			Usage:  "send message to user",
			EnvVar: "PLUGIN_TO",
		},
		cli.StringSliceFlag{
			Name:   "message",
			Usage:  "line message",
			EnvVar: "PLUGIN_MESSAGE",
		},
		cli.StringSliceFlag{
			Name:   "image",
			Usage:  "line image",
			EnvVar: "PLUGIN_IMAGE",
		},
		cli.StringSliceFlag{
			Name:   "video",
			Usage:  "line video",
			EnvVar: "PLUGIN_VIDEO",
		},
		cli.StringSliceFlag{
			Name:   "audio",
			Usage:  "line audio",
			EnvVar: "PLUGIN_AUDIO",
		},
		cli.StringSliceFlag{
			Name:   "sticker",
			Usage:  "line sticker",
			EnvVar: "PLUGIN_STICKER",
		},
		cli.StringSliceFlag{
			Name:   "location",
			Usage:  "line location",
			EnvVar: "PLUGIN_LOCATION",
		},
		cli.StringFlag{
			Name:   "delimiter",
			Usage:  "line delimiter",
			Value:  "::",
			EnvVar: "PLUGIN_DELIMITER",
		},
		cli.StringFlag{
			Name:   "repo.owner",
			Usage:  "repository owner",
			EnvVar: "DRONE_REPO_OWNER",
		},
		cli.StringFlag{
			Name:   "repo.name",
			Usage:  "repository name",
			EnvVar: "DRONE_REPO_NAME",
		},
		cli.StringFlag{
			Name:   "commit.sha",
			Usage:  "git commit sha",
			EnvVar: "DRONE_COMMIT_SHA",
		},
		cli.StringFlag{
			Name:   "commit.branch",
			Value:  "master",
			Usage:  "git commit branch",
			EnvVar: "DRONE_COMMIT_BRANCH",
		},
		cli.StringFlag{
			Name:   "commit.author",
			Usage:  "git author name",
			EnvVar: "DRONE_COMMIT_AUTHOR",
		},
		cli.StringFlag{
			Name:   "build.event",
			Value:  "push",
			Usage:  "build event",
			EnvVar: "DRONE_BUILD_EVENT",
		},
		cli.IntFlag{
			Name:   "build.number",
			Usage:  "build number",
			EnvVar: "DRONE_BUILD_NUMBER",
		},
		cli.StringFlag{
			Name:   "build.status",
			Usage:  "build status",
			Value:  "success",
			EnvVar: "DRONE_BUILD_STATUS",
		},
		cli.StringFlag{
			Name:   "build.link",
			Usage:  "build link",
			EnvVar: "DRONE_BUILD_LINK",
		},
	}
	app.Run(os.Args)
}

func run(c *cli.Context) error {
	plugin := Plugin{
		Repo: Repo{
			Owner: c.String("repo.owner"),
			Name:  c.String("repo.name"),
		},
		Build: Build{
			Number: c.Int("build.number"),
			Event:  c.String("build.event"),
			Status: c.String("build.status"),
			Commit: c.String("commit.sha"),
			Branch: c.String("commit.branch"),
			Author: c.String("commit.author"),
			Link:   c.String("build.link"),
		},
		Config: Config{
			ChannelID:     c.String("channel.id"),
			ChannelSecret: c.String("channel.secret"),
			MID:           c.String("mid"),
			To:            c.StringSlice("to"),
			Delimiter:     c.String("delimiter"),
			Message:       c.StringSlice("message"),
			Image:         c.StringSlice("image"),
			Video:         c.StringSlice("video"),
			Audio:         c.StringSlice("audio"),
			Sticker:       c.StringSlice("sticker"),
			Location:      c.StringSlice("location"),
		},
	}

	return plugin.Exec()
}
