package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	"github.com/urfave/cli"
)

// Version set at compile-time
var (
	version = "0.0.0"
	build   = "0"
)

func main() {
	app := cli.NewApp()
	app.Name = "Drone LINE"
	app.Usage = "Send LINE notification"
	app.Copyright = "Copyright (c) 2018 Bo-Yi Wu"
	app.Version = fmt.Sprintf("%s+%s", version, build)
	app.Authors = []cli.Author{
		{
			Name:  "Bo-Yi Wu",
			Email: "appleboy.tw@gmail.com",
		},
	}
	app.Action = run
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "secret, s",
			Usage:  "line channel secret",
			EnvVar: "PLUGIN_CHANNEL_SECRET,LINE_CHANNEL_SECRET,INPUT_SECRET",
		},
		cli.StringFlag{
			Name:   "token, t",
			Usage:  "line channel access token",
			EnvVar: "PLUGIN_CHANNEL_TOKEN,LINE_CHANNEL_TOKEN,INPUT_TOKEN",
		},
		cli.StringSliceFlag{
			Name:   "to, u",
			Usage:  "line user ID",
			EnvVar: "PLUGIN_TO,LINE_TO,INPUT_TO",
		},
		cli.StringFlag{
			Name:   "toroom, r",
			Usage:  "line room ID",
			EnvVar: "PLUGIN_TO_ROOM,LINE_TO_ROOM,INPUT_ROOM",
		},
		cli.StringFlag{
			Name:   "togroup, g",
			Usage:  "line group ID",
			EnvVar: "PLUGIN_TO_GROUP,LINE_TO_GROUP,INPUT_GROUP",
		},
		cli.StringSliceFlag{
			Name:   "message, m",
			Usage:  "line message",
			EnvVar: "PLUGIN_MESSAGE,LINE_MESSAGE,INPUT_MESSAGE",
		},
		cli.StringSliceFlag{
			Name:   "image",
			Usage:  "line image",
			EnvVar: "PLUGIN_IMAGES,LINE_IMAGES,INPUT_IMAGES",
		},
		cli.StringSliceFlag{
			Name:   "video",
			Usage:  "line video",
			EnvVar: "PLUGIN_VIDEOS,LINE_VIDEOS,INPUT_VIDEOS",
		},
		cli.StringSliceFlag{
			Name:   "audio",
			Usage:  "line audio",
			EnvVar: "PLUGIN_AUDIOS,LINE_AUDIOS,INPUT_AUDIOS",
		},
		cli.StringSliceFlag{
			Name:   "sticker",
			Usage:  "line sticker",
			EnvVar: "PLUGIN_STICKERS,LINE_STICKERS,INPUT_STICKERS",
		},
		cli.StringSliceFlag{
			Name:   "location",
			Usage:  "line location",
			EnvVar: "PLUGIN_LOCATIONS,LINE_LOCATIONS,INPUT_LOCATIONS",
		},
		cli.StringFlag{
			Name:   "delimiter",
			Usage:  "line delimiter",
			Value:  "::",
			EnvVar: "PLUGIN_DELIMITER,LINE_DELIMITER,INPUT_DELIMITER",
		},
		cli.BoolFlag{
			Name:   "match.email",
			Usage:  "send message when only match email",
			EnvVar: "PLUGIN_ONLY_MATCH_EMAIL",
		},
		cli.IntFlag{
			Name:   "port, P",
			Usage:  "webhook port",
			EnvVar: "LINE_PORT",
			Value:  8088,
		},
		cli.BoolFlag{
			Name:   "drone",
			Usage:  "environment is drone",
			EnvVar: "DRONE",
		},
		cli.StringFlag{
			Name:   "repo",
			Usage:  "repository owner and repository name",
			EnvVar: "DRONE_REPO,GITHUB_REPOSITORY",
		},
		cli.StringFlag{
			Name:   "repo.namespace",
			Usage:  "repository namespace",
			EnvVar: "DRONE_REPO_OWNER,DRONE_REPO_NAMESPACE,GITHUB_ACTOR",
		},
		cli.StringFlag{
			Name:   "repo.name",
			Usage:  "repository name",
			EnvVar: "DRONE_REPO_NAME",
		},
		cli.StringFlag{
			Name:   "commit.sha",
			Usage:  "git commit sha",
			EnvVar: "DRONE_COMMIT_SHA,GITHUB_SHA",
		},
		cli.StringFlag{
			Name:   "commit.ref",
			Usage:  "git commit ref",
			EnvVar: "DRONE_COMMIT_REF,GITHUB_REF",
		},
		cli.StringFlag{
			Name:   "commit.branch",
			Value:  "master",
			Usage:  "git commit branch",
			EnvVar: "DRONE_COMMIT_BRANCH",
		},
		cli.StringFlag{
			Name:   "commit.link",
			Usage:  "git commit link",
			EnvVar: "DRONE_COMMIT_LINK",
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
			Name:   "commit.author.avatar",
			Usage:  "git author avatar",
			EnvVar: "DRONE_COMMIT_AUTHOR_AVATAR",
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
		cli.StringFlag{
			Name:   "pull.request",
			Usage:  "pull request",
			EnvVar: "DRONE_PULL_REQUEST",
		},
		cli.Float64Flag{
			Name:   "job.started",
			Usage:  "job started",
			EnvVar: "DRONE_BUILD_STARTED",
		},
		cli.Float64Flag{
			Name:   "job.finished",
			Usage:  "job finished",
			EnvVar: "DRONE_BUILD_FINISHED",
		},
		cli.StringFlag{
			Name:  "env-file",
			Usage: "source env file",
		},
		cli.BoolFlag{
			Name:   "github",
			Usage:  "Boolean value, indicates the runtime environment is GitHub Action.",
			EnvVar: "PLUGIN_GITHUB,GITHUB",
		},
		cli.StringFlag{
			Name:   "github.workflow",
			Usage:  "The name of the workflow.",
			EnvVar: "GITHUB_WORKFLOW",
		},
		cli.StringFlag{
			Name:   "github.action",
			Usage:  "The name of the action.",
			EnvVar: "GITHUB_ACTION",
		},
		cli.StringFlag{
			Name:   "github.event.name",
			Usage:  "The webhook name of the event that triggered the workflow.",
			EnvVar: "GITHUB_EVENT_NAME",
		},
		cli.StringFlag{
			Name:   "github.event.path",
			Usage:  "The path to a file that contains the payload of the event that triggered the workflow. Value: /github/workflow/event.json.",
			EnvVar: "GITHUB_EVENT_PATH",
		},
		cli.StringFlag{
			Name:   "github.workspace",
			Usage:  "The GitHub workspace path. Value: /github/workspace.",
			EnvVar: "GITHUB_WORKSPACE",
		},
		cli.StringFlag{
			Name:   "deploy.to",
			Usage:  "Provides the target deployment environment for the running build. This value is only available to promotion and rollback pipelines.",
			EnvVar: "DRONE_DEPLOY_TO",
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
			EnvVar: "PLUGIN_DOMAIN,DOMAIN",
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

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	if c.String("env-file") != "" {
		_ = godotenv.Load(c.String("env-file"))
	}

	plugin := Plugin{
		GitHub: GitHub{
			Workflow:  c.String("github.workflow"),
			Workspace: c.String("github.workspace"),
			Action:    c.String("github.action"),
			EventName: c.String("github.event.name"),
			EventPath: c.String("github.event.path"),
		},
		Repo: Repo{
			FullName:  c.String("repo"),
			Namespace: c.String("repo.namespace"),
			Name:      c.String("repo.name"),
		},
		Commit: Commit{
			Sha:     c.String("commit.sha"),
			Ref:     c.String("commit.ref"),
			Branch:  c.String("commit.branch"),
			Link:    c.String("commit.link"),
			Author:  c.String("commit.author"),
			Email:   c.String("commit.author.email"),
			Avatar:  c.String("commit.author.avatar"),
			Message: c.String("commit.message"),
		},
		Build: Build{
			Tag:      c.String("build.tag"),
			Number:   c.Int("build.number"),
			Event:    c.String("build.event"),
			Status:   c.String("build.status"),
			Link:     c.String("build.link"),
			Started:  c.Float64("job.started"),
			Finished: c.Float64("job.finished"),
			PR:       c.String("pull.request"),
			DeployTo: c.String("deploy.to"),
		},
		Config: Config{
			ChannelSecret: c.String("secret"),
			ChannelToken:  c.String("token"),
			To:            c.StringSlice("to"),
			ToRoom:        c.String("toroom"),
			ToGroup:       c.String("togroup"),
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
			AutoTLS:       c.Bool("autotls"),
			Host:          c.StringSlice("host"),
		},
	}

	command := c.Args().Get(0)

	switch command {
	case "webhook":
		return plugin.Webhook()
	case "notify":
		return plugin.Notify()
	}

	return plugin.Exec()
}
