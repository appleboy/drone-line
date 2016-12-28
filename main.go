package main

import (
	"os"

	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	"github.com/urfave/cli"
)

// Version set at compile-time
var Version string

func main() {
	app := cli.NewApp()
	app.Name = "line plugin"
	app.Usage = "line plugin"
	app.Action = run
	app.Version = Version
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "secret",
			Usage:  "line channel secret",
			EnvVar: "PLUGIN_CHANNEL_SECRET,LINE_CHANNEL_SECRET",
		},
		cli.StringFlag{
			Name:   "token",
			Usage:  "line channel access token",
			EnvVar: "PLUGIN_CHANNEL_TOKEN,LINE_CHANNEL_TOKEN",
		},
		cli.StringSliceFlag{
			Name:   "to",
			Usage:  "line user ID",
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
		cli.BoolFlag{
			Name:   "match.email",
			Usage:  "send message when only match email",
			EnvVar: "PLUGIN_ONLY_MATCH_EMAIL",
		},
		cli.IntFlag{
			Name:   "port",
			Usage:  "webhook port",
			EnvVar: "PLUGIN_PORT",
			Value:  8088,
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
			Name:   "commit.author.email",
			Usage:  "git author email",
			EnvVar: "DRONE_COMMIT_AUTHOR_EMAIL",
		},
		cli.StringFlag{
			Name:   "commit.message",
			Usage:  "commit message",
			EnvVar: "DRONE_COMMIT_MESSAGE",
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
		cli.StringFlag{
			Name:   "build.tag",
			Usage:  "build tag",
			EnvVar: "DRONE_TAG",
		},
		cli.Float64Flag{
			Name:   "job.started",
			Usage:  "job started",
			EnvVar: "DRONE_JOB_STARTED",
		},
		cli.Float64Flag{
			Name:   "job.finished",
			Usage:  "job finished",
			EnvVar: "DRONE_JOB_FINISHED",
		},
		cli.StringFlag{
			Name:  "env-file",
			Usage: "source env file",
		},
	}
	app.Run(os.Args)
}

func run(c *cli.Context) error {
	if c.String("env-file") != "" {
		_ = godotenv.Load(c.String("env-file"))
	}

	plugin := Plugin{
		Repo: Repo{
			Owner: c.String("repo.owner"),
			Name:  c.String("repo.name"),
		},
		Build: Build{
			Tag:      c.String("build.tag"),
			Number:   c.Int("build.number"),
			Event:    c.String("build.event"),
			Status:   c.String("build.status"),
			Commit:   c.String("commit.sha"),
			Branch:   c.String("commit.branch"),
			Author:   c.String("commit.author"),
			Email:    c.String("commit.author.email"),
			Message:  c.String("commit.message"),
			Link:     c.String("build.link"),
			Started:  c.Float64("job.started"),
			Finished: c.Float64("job.finished"),
		},
		Config: Config{
			ChannelSecret: c.String("secret"),
			ChannelToken:  c.String("token"),
			To:            c.StringSlice("to"),
			Delimiter:     c.String("delimiter"),
			MatchEmail:    c.Bool("match.email"),
			Message:       c.StringSlice("message"),
			Image:         c.StringSlice("image"),
			Video:         c.StringSlice("video"),
			Audio:         c.StringSlice("audio"),
			Sticker:       c.StringSlice("sticker"),
			Location:      c.StringSlice("location"),
			Port:          c.Int("port"),
		},
	}

	command := c.Args().Get(0)

	if command == "webhook" {
		return plugin.Webhook()
	}

	return plugin.Exec()
}
