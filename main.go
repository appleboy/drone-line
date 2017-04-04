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
	app.Name = "Drone LINE"
	app.Usage = "Send LINE notification"
	app.Copyright = "Copyright (c) 2017 Bo-Yi Wu"
	app.Authors = []cli.Author{
		{
			Name:  "Bo-Yi Wu",
			Email: "appleboy.tw@gmail.com",
		},
	}
	app.Action = run
	app.Version = Version
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "secret, s",
			Usage:  "line channel secret",
			EnvVar: "PLUGIN_CHANNEL_SECRET,LINE_CHANNEL_SECRET",
		},
		cli.StringFlag{
			Name:   "token, t",
			Usage:  "line channel access token",
			EnvVar: "PLUGIN_CHANNEL_TOKEN,LINE_CHANNEL_TOKEN",
		},
		cli.StringSliceFlag{
			Name:   "to, u",
			Usage:  "line user ID",
			EnvVar: "PLUGIN_TO,LINE_TO",
		},
		cli.StringSliceFlag{
			Name:   "message, m",
			Usage:  "line message",
			EnvVar: "PLUGIN_MESSAGE,LINE_MESSAGE",
		},
		cli.StringSliceFlag{
			Name:   "image",
			Usage:  "line image",
			EnvVar: "PLUGIN_IMAGES,LINE_IMAGES",
		},
		cli.StringSliceFlag{
			Name:   "video",
			Usage:  "line video",
			EnvVar: "PLUGIN_VIDEOS,LINE_VIDEOS",
		},
		cli.StringSliceFlag{
			Name:   "audio",
			Usage:  "line audio",
			EnvVar: "PLUGIN_AUDIOS,LINE_AUDIOS",
		},
		cli.StringSliceFlag{
			Name:   "sticker",
			Usage:  "line sticker",
			EnvVar: "PLUGIN_STICKERS,LINE_STICKERS",
		},
		cli.StringSliceFlag{
			Name:   "location",
			Usage:  "line location",
			EnvVar: "PLUGIN_LOCATIONS,LINE_LOCATIONS",
		},
		cli.StringFlag{
			Name:   "delimiter",
			Usage:  "line delimiter",
			Value:  "::",
			EnvVar: "PLUGIN_DELIMITER,LINE_DELIMITER",
		},
		cli.BoolFlag{
			Name:   "match.email",
			Usage:  "send message when only match email",
			EnvVar: "PLUGIN_ONLY_MATCH_EMAIL",
		},
		cli.IntFlag{
			Name:   "port, P",
			Usage:  "webhook port",
			EnvVar: "LINE_WEBHOOK_PORT",
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
		cli.BoolFlag{
			Name:   "tunnel",
			Usage:  "Enable tunnel host for webhook",
			EnvVar: "PLUGIN_TUNNEL,TUNNEL",
		},
		cli.BoolFlag{
			Name:   "debug",
			Usage:  "debug mode",
			EnvVar: "PLUGIN_DEBUG,DEBUG",
		},
		cli.StringFlag{
			Name:   "domain",
			Usage:  "tunnel host name must be lowercase and between 4 and 63 alphanumeric characters.",
			EnvVar: "DRONE_JOB_FINISHED",
		},
		cli.BoolFlag{
			Name:   "autotls",
			Usage:  "Auto tls mode",
			EnvVar: "PLUGIN_AUTOTLS,AUTOTLS",
		},
		cli.StringSliceFlag{
			Name:   "host",
			Usage:  "Auto tls host name",
			EnvVar: "PLUGIN_HOSTNAME,HOSTNAME",
		},
		cli.StringFlag{
			Name:   "cache",
			Usage:  "folder for storing certificates",
			EnvVar: "PLUGIN_CACHE,CACHE",
		},
	}

	// Override a template
	cli.AppHelpTemplate = `
________                                       .____    .___ _______  ___________
\______ \_______  ____   ____   ____           |    |   |   |\      \ \_   _____/
 |    |  \_  __ \/  _ \ /    \_/ __ \   ______ |    |   |   |/   |   \ |    __)_
 |    |   \  | \(  <_> )   |  \  ___/  /_____/ |    |___|   /    |    \|        \
/_______  /__|   \____/|___|  /\___  >         |_______ \___\____|__  /_______  /
        \/                  \/     \/                  \/           \/        \/
                                                              version: {{.Version}}
NAME:
   {{.Name}} - {{.Usage}}

USAGE:
   {{.HelpName}} {{if .VisibleFlags}}[global options]{{end}}{{if .Commands}} command [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}
   {{if len .Authors}}
AUTHOR:
   {{range .Authors}}{{ . }}{{end}}
   {{end}}{{if .Commands}}
COMMANDS:
{{range .Commands}}{{if not .HideHelp}}   {{join .Names ", "}}{{ "\t"}}{{.Usage}}{{ "\n" }}{{end}}{{end}}{{end}}{{if .VisibleFlags}}
GLOBAL OPTIONS:
   {{range .VisibleFlags}}{{.}}
   {{end}}{{end}}{{if .Copyright }}
COPYRIGHT:
   {{.Copyright}}
   {{end}}{{if .Version}}
VERSION:
   {{.Version}}
   {{end}}
REPOSITORY:
    Github: https://github.com/appleboy/drone-line
`
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
			Tunnel:        c.Bool("tunnel"),
			Debug:         c.Bool("debug"),
			Domain:        c.String("domain"),
			AutoTLS:       c.Bool("AutoTLS"),
			Host:          c.StringSlice("host"),
			Cache:         c.String("cache"),
		},
	}

	command := c.Args().Get(0)

	if command == "webhook" {
		return plugin.Webhook()
	}

	return plugin.Exec()
}
