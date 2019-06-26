package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/ebracho/kanjibot/jisho"
	"github.com/ebracho/twitchbot"
)

const (
	defaultConfigPath = "config.json"
)

var cfg struct {
	TwitchUser  string   `json:"twitchUser"`
	TwitchToken string   `json:"twitchToken"`
	Channels    []string `json:"channels"`
}

var (
	reKanjiSearch = regexp.MustCompile(`!k ([a-zA-Z ]|\p{Han}|\p{Hiragana}|\p{Katakana})+`)
)

func main() {
	cfgPath := flag.String("configPath", defaultConfigPath, "path to kanjibot config file")
	flag.Parse()

	// Parse config file
	cfgFile, err := os.Open(*cfgPath)
	if err != nil {
		log.Fatal(err)
	}
	d := json.NewDecoder(cfgFile)
	if err := d.Decode(&cfg); err != nil {
		log.Fatal(err)
	}

	j := jisho.New()
	client := twitchbot.New(cfg.TwitchUser, cfg.Channels, cfg.TwitchToken)
	client.RegisterHandler(reKanjiSearch, func(channel, text string, c *twitchbot.Client) {
		parts := strings.Split(text, " ")
		if len(parts) < 2 {
			return
		}
		keyword := strings.Join(parts[1:], " ")
		res, err := j.SearchWords(keyword)
		if err != nil {
			log.Println(err)
			return
		}
		if len(res.Data) == 0 {
			return
		}
		d := res.Data[0]
		if len(d.Japanese) == 0 {
			return
		}
		jp := d.Japanese[0]
		msg := fmt.Sprintf("%s (%s)", jp.Word, jp.Reading)
		if len(d.Senses) > 0 {
			s := d.Senses[0]
			if len(s.EnglishDefinitions) > 0 {
				msg += ": " + strings.Join(s.EnglishDefinitions, ", ")
			}
			if len(s.PartsOfSpeech) > 0 {
				msg += fmt.Sprintf(" (%s)", strings.Join(s.PartsOfSpeech, ", "))
			}
		}
		c.MessageChannel(channel, msg)
	})
	if err := client.Run(); err != nil {
		log.Fatal(err.Error())
	}
}
