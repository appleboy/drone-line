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
	app.Name = "facebook plugin"
	app.Usage = "facebook plugin"
	app.Action = run
	app.Version = Version
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "page.token",
			Usage:  "facebook page token",
			EnvVar: "PLUGIN_FB_PAGE_TOKEN,FB_PAGE_TOKEN,PAGE_TOKEN",
		},
		cli.StringFlag{
			Name:   "verify.token",
			Usage:  "facebook verify token",
			EnvVar: "PLUGIN_FB_VERIFY_TOKEN,FB_VERIFY_TOKEN,VERIFY_TOKEN",
		},
		cli.BoolFlag{
			Name:   "verify",
			Usage:  "verifying webhooks on the Facebook Developer Portal",
			EnvVar: "PLUGIN_VERIFY,VERIFY",
		},
		cli.StringSliceFlag{
			Name:   "to",
			Usage:  "send message to user",
			EnvVar: "PLUGIN_TO",
		},
		cli.StringSliceFlag{
			Name:   "message",
			Usage:  "text message",
			EnvVar: "PLUGIN_MESSAGE",
		},
		cli.StringSliceFlag{
			Name:   "image",
			Usage:  "image message",
			EnvVar: "PLUGIN_IMAGES",
		},
		cli.StringSliceFlag{
			Name:   "audio",
			Usage:  "audio message",
			EnvVar: "PLUGIN_AUDIOS",
		},
		cli.StringSliceFlag{
			Name:   "video",
			Usage:  "video message",
			EnvVar: "PLUGIN_VIDEOS",
		},
		cli.StringSliceFlag{
			Name:   "file",
			Usage:  "file message",
			EnvVar: "PLUGIN_FILES",
		},
		cli.BoolFlag{
			Name:   "match.email",
			Usage:  "send message when only match email",
			EnvVar: "PLUGIN_ONLY_MATCH_EMAIL",
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
			Name:   "env-file",
			Usage:  "source env file",
			EnvVar: "ENV_FILE",
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
			PageToken:   c.String("page.token"),
			VerifyToken: c.String("verify.token"),
			Verify:      c.Bool("verify"),
			MatchEmail:  c.Bool("match.email"),
			To:          c.StringSlice("to"),
			Message:     c.StringSlice("message"),
			Image:       c.StringSlice("image"),
			Audio:       c.StringSlice("audio"),
			Video:       c.StringSlice("video"),
			File:        c.StringSlice("file"),
		},
	}

	return plugin.Exec()
}
