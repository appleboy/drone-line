package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/NoahShen/gotunnelme/src/gotunnelme"
	"github.com/appleboy/com/random"
	"github.com/drone/drone-template-lib/template"
	"github.com/fatih/color"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/crypto/acme/autocert"
)

const defaultPreviewImageURL = "https://cdn4.iconfinder.com/data/icons/miu/24/device-camera-recorder-video-glyph-256.png"

type (
	// GitHub information.
	GitHub struct {
		Workflow  string
		Workspace string
		Action    string
		EventName string
		EventPath string
	}

	// Repo information.
	Repo struct {
		FullName  string
		Namespace string
		Name      string
	}

	// Commit information.
	Commit struct {
		Sha     string
		Ref     string
		Branch  string
		Link    string
		Author  string
		Avatar  string
		Email   string
		Message string
	}

	// Build information.
	Build struct {
		Tag      string
		Event    string
		Number   int
		Status   string
		Link     string
		Started  float64
		Finished float64
		PR       string
		DeployTo string
	}

	// Config for the plugin.
	Config struct {
		ChannelToken  string
		ChannelSecret string
		To            []string
		ToRoom        string
		ToGroup       string
		Delimiter     string
		Message       []string
		Image         []string
		Video         []string
		Audio         []string
		Sticker       []string
		Location      []string
		MatchEmail    bool
		Port          int
		Tunnel        bool
		Debug         bool
		Domain        string
		AutoTLS       bool
		Host          []string
	}

	// Plugin values.
	Plugin struct {
		Repo   Repo
		Build  Build
		Config Config
		Commit Commit
		GitHub GitHub
	}

	// Audio format
	Audio struct {
		URL      string
		Duration int
	}

	// Location format
	Location struct {
		Title     string
		Address   string
		Latitude  float64
		Longitude float64
	}
)

var (
	// ReceiveCount is receive notification count
	ReceiveCount int64
	// SendCount is send notification count
	SendCount int64
)

func init() {
	// Support metrics
	m := NewMetrics()
	prometheus.MustRegister(m)
}

func trimElement(keys []string) []string {
	var newKeys []string

	for _, value := range keys {
		value = strings.Trim(value, " ")
		if len(value) == 0 {
			continue
		}
		newKeys = append(newKeys, value)
	}

	return newKeys
}

func convertImage(value, delimiter string) []string {
	values := trimElement(strings.Split(value, delimiter))

	if len(values) < 2 {
		values = append(values, values[0])
	}

	return values
}

func convertVideo(value, delimiter string) []string {
	values := trimElement(strings.Split(value, delimiter))

	if len(values) < 2 {
		values = append(values, defaultPreviewImageURL)
	}

	return values
}

func convertAudio(value, delimiter string) (Audio, bool) {
	values := trimElement(strings.Split(value, delimiter))

	if len(values) < 2 {
		return Audio{}, true
	}

	duration, err := strconv.Atoi(values[1])

	if err != nil {
		log.Println(err.Error())
		return Audio{}, true
	}

	return Audio{
		URL:      values[0],
		Duration: duration,
	}, false
}

func convertSticker(value, delimiter string) ([]string, bool) {
	values := trimElement(strings.Split(value, delimiter))

	if len(values) < 2 {
		return []string{}, true
	}

	return values, false
}

func convertLocation(value, delimiter string) (Location, bool) {
	var latitude, longitude float64
	var err error
	values := trimElement(strings.Split(value, delimiter))

	if len(values) < 4 {
		return Location{}, true
	}

	latitude, err = strconv.ParseFloat(values[2], 64)

	if err != nil {
		log.Println(err.Error())
		return Location{}, true
	}

	longitude, err = strconv.ParseFloat(values[3], 64)

	if err != nil {
		log.Println(err.Error())
		return Location{}, true
	}

	return Location{
		Title:     values[0],
		Address:   values[1],
		Latitude:  latitude,
		Longitude: longitude,
	}, false
}

func parseTo(to []string, authorEmail string, matchEmail bool, delimiter string) []string {
	var emails []string
	var ids []string
	attachEmail := true

	for _, value := range trimElement(to) {
		idArray := trimElement(strings.Split(value, delimiter))

		// check match author email
		if len(idArray) > 1 {
			if email := idArray[1]; email != authorEmail {
				continue
			}

			emails = append(emails, idArray[0])
			attachEmail = false
			continue
		}

		ids = append(ids, idArray[0])
	}

	if matchEmail == true && attachEmail == false {
		return emails
	}

	for _, value := range emails {
		ids = append(ids, value)
	}

	return ids
}

// Bot is new Line Bot clien.
func (p Plugin) Bot() (*linebot.Client, error) {
	if len(p.Config.ChannelToken) == 0 || len(p.Config.ChannelSecret) == 0 {
		return nil, errors.New("missing line bot config")
	}

	return linebot.New(p.Config.ChannelSecret, p.Config.ChannelToken)
}

func (p Plugin) getTunnelDomain() (string, error) {
	var domain string
	if p.Config.Domain != "" {
		if len(p.Config.Domain) < 4 || len(p.Config.Domain) > 63 {
			return "", errors.New("tunnel host name must be lowercase and between 4 and 63 alphanumeric characters")
		}
		domain = p.Config.Domain
	} else {
		domain = strings.ToLower(random.String(10))
	}

	return domain, nil
}

// Handler is http handler.
func (p Plugin) Handler(bot *linebot.Client) *http.ServeMux {
	mux := http.NewServeMux()

	// Setup HTTP Server for receiving requests from LINE platform
	mux.HandleFunc("/callback", func(w http.ResponseWriter, req *http.Request) {
		events, err := bot.ParseRequest(req)
		if err != nil {
			if err == linebot.ErrInvalidSignature {
				w.WriteHeader(400)
			} else {
				w.WriteHeader(500)
			}
			return
		}
		for _, event := range events {
			if event.Type == linebot.EventTypeMessage {
				switch message := event.Message.(type) {
				case *linebot.TextMessage:
					log.Printf("%#v", event.Source)
					log.Printf("User ID is %v\n", event.Source.UserID)
					log.Printf("Room ID is %v\n", event.Source.RoomID)
					log.Printf("Group ID is %v\n", event.Source.GroupID)

					ReceiveCount++
					if message.Text == "test" {
						SendCount++
						if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("count + 1")).Do(); err != nil {
							log.Print(err)
						}
					}
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message.Text)).Do(); err != nil {
						log.Print(err)
					}
				}
			}
		}
	})

	mux.HandleFunc("/metrics", func(w http.ResponseWriter, req *http.Request) {
		promhttp.Handler().ServeHTTP(w, req)
	})

	// Setup HTTP Server for receiving requests from LINE platform
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(w, "Welcome to Line webhook page.")
	})

	return mux
}

// Webhook support line callback service.
func (p Plugin) Webhook() error {
	readyToListen := false
	bot, err := p.Bot()

	if err != nil {
		return err
	}
	mux := p.Handler(bot)

	if p.Config.Tunnel {
		if p.Config.Debug {
			gotunnelme.Debug = true
		}

		domain, err := p.getTunnelDomain()
		if err != nil {
			panic(err)
		}

		tunnel := gotunnelme.NewTunnel()
		url, err := tunnel.GetUrl(domain)
		if err != nil {
			panic("Could not get localtunnel.me URL. " + err.Error())
		}
		go func() {
			for !readyToListen {
				time.Sleep(1 * time.Second)
			}
			c := color.New(color.FgYellow)
			c.Println("Tunnel URL:", url)
			err := tunnel.CreateTunnel(p.Config.Port)
			if err != nil {
				panic("Could not create tunnel. " + err.Error())
			}
		}()
	}

	readyToListen = true
	if p.Config.Port != 443 && !p.Config.AutoTLS {
		log.Println("Line Webhook Server Listin on " + strconv.Itoa(p.Config.Port) + " port")
		if err := http.ListenAndServe(":"+strconv.Itoa(p.Config.Port), mux); err != nil {
			log.Fatal(err)
		}
	}

	if p.Config.AutoTLS && len(p.Config.Host) != 0 {
		log.Println("Line Webhook Server Listin on 443 port, hostname: " + strings.Join(p.Config.Host, ", "))
		return http.Serve(autocert.NewListener(p.Config.Host...), mux)
	}

	return nil
}

// Notify for Line notify service
// https://notify-bot.line.me
func (p Plugin) Notify() error {
	if p.Config.ChannelToken == "" || len(p.Config.Message) == 0 {
		return errors.New("missing token or message")
	}

	for _, m := range p.Config.Message {
		if err := p.notify(m, p.Config.ChannelToken); err != nil {
			return err
		}
	}

	return nil
}

func (p Plugin) notify(message, token string) error {
	data := url.Values{}
	data.Add("message", message)

	u, _ := url.ParseRequestURI("https://notify-api.line.me/api/notify")
	urlStr := u.String()

	req, err := http.NewRequest(
		"POST",
		urlStr,
		strings.NewReader(data.Encode()),
	)

	if err != nil {
		return errors.New("failed to create request:" + err.Error())
	}

	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		return errors.New("failed to process request:" + err.Error())
	}

	defer res.Body.Close()

	if p.Config.Debug {
		log.Println("=================================")
		log.Printf("%#v\n", res)
		log.Println("=================================")
	}

	if res.Status == "200 OK" {
		log.Println("successfully send notfiy")
	}

	if err != nil {
		return errors.New("failed to create request:" + err.Error())
	}

	return nil
}

// Exec executes the plugin.
func (p Plugin) Exec() error {

	bot, err := p.Bot()

	if err != nil {
		return err
	}

	if len(p.Config.To) == 0 && len(p.Config.ToRoom) == 0 && len(p.Config.ToGroup) == 0 {
		return errors.New("missing line user config")
	}

	var message []string
	if len(p.Config.Message) > 0 {
		message = p.Config.Message
	} else {
		message = p.Message(p.Repo, p.Build, p.Commit)
	}

	// Initial messages array.
	var messages []linebot.SendingMessage

	for _, value := range trimElement(message) {
		txt, err := template.RenderTrim(value, p)
		if err != nil {
			return err
		}

		messages = append(messages, linebot.NewTextMessage(txt))
	}

	// Add image message
	for _, value := range trimElement(p.Config.Image) {
		values := convertImage(value, p.Config.Delimiter)

		messages = append(messages, linebot.NewImageMessage(values[0], values[1]))
	}

	// Add image message.
	for _, value := range trimElement(p.Config.Video) {
		values := convertVideo(value, p.Config.Delimiter)

		messages = append(messages, linebot.NewVideoMessage(values[0], values[1]))
	}

	// Add Audio message.
	for _, value := range trimElement(p.Config.Audio) {
		audio, empty := convertAudio(value, p.Config.Delimiter)

		if empty == true {
			continue
		}

		messages = append(messages, linebot.NewAudioMessage(audio.URL, audio.Duration))
	}

	// Add Sticker message.
	for _, value := range trimElement(p.Config.Sticker) {
		sticker, empty := convertSticker(value, p.Config.Delimiter)

		if empty == true {
			continue
		}

		messages = append(messages, linebot.NewStickerMessage(sticker[0], sticker[1]))
	}

	// check Location array.
	for _, value := range trimElement(p.Config.Location) {
		location, empty := convertLocation(value, p.Config.Delimiter)

		if empty == true {
			continue
		}

		messages = append(messages, linebot.NewLocationMessage(location.Title, location.Address, location.Latitude, location.Longitude))
	}

	uids := parseTo(p.Config.To, p.Commit.Email, p.Config.MatchEmail, p.Config.Delimiter)

	// Send messages to multiple users at any time.
	if len(uids) > 0 {
		if _, err := bot.Multicast(uids, messages...).Do(); err != nil {
			log.Println(err.Error())
		}
	}

	// Send messages to single room at any time.
	rid := p.Config.ToRoom
	if rid != "" {
		if _, err := bot.PushMessage(rid, messages...).Do(); err != nil {
			log.Println(err.Error())
		}
	}

	// Send messages to single group at any time.
	gid := p.Config.ToGroup
	if gid != "" {
		if _, err := bot.PushMessage(gid, messages...).Do(); err != nil {
			log.Println(err.Error())
		}
	}

	return nil
}

// Message is line default message.
func (p Plugin) Message(repo Repo, build Build, commit Commit) []string {
	return []string{fmt.Sprintf("[%s] <%s> (%s)『%s』by %s",
		build.Status,
		build.Link,
		commit.Branch,
		commit.Message,
		commit.Author,
	)}
}
