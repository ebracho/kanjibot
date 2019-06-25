package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strings"

	"github.com/ebracho/kanjibot/jisho"
	"github.com/ebracho/twitchbot"
)

type stringSliceFlags []string

func (s *stringSliceFlags) String() string {
	return strings.Join(*s, ",")
}

func (s *stringSliceFlags) Set(value string) error {
	*s = append(*s, value)
	return nil
}

var (
	twitchUser      string
	twitchTokenPath string
	twitchChannels  stringSliceFlags
)

var (
	reKanjiSearch = regexp.MustCompile(`!k ([a-zA-Z ]|\p{Han}|\p{Hiragana}|\p{Katakana})+`)
)

func main() {
	flag.StringVar(&twitchUser, "twitch_user", "", "twitch user name")
	flag.StringVar(&twitchTokenPath, "twitch_token_file", "", "twitch oauth token filepath")
	flag.Var(&twitchChannels, "channel", "twitch channel to join")
	flag.Parse()

	if twitchUser == "" || twitchTokenPath == "" || twitchChannels == nil {
		log.Fatalf("Must specify -twitch_user, -twitch_token_file, -channel")
	}

	tokenBytes, err := ioutil.ReadFile(twitchTokenPath)
	if err != nil {
		log.Fatal(err)
	}
	token := strings.TrimSpace(string(tokenBytes))

	j := jisho.New()
	client := twitchbot.New(twitchUser, twitchChannels, token)
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
