package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"strings"

	"github.com/bwmarrin/discordgo"
	lyrics "github.com/rhnvrm/lyric-api-go"
)

var (
	Token     string
	BotPrefix string

	config *configStruct
)

type configStruct struct {
	Token     string `json : "Token"`
	BotPrefix string `json : "BotPrefix"`
}

func ReadConfig() error {
	fmt.Println("Reading config file...")
	file, err := ioutil.ReadFile("./config.json")

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println(string(file))

	err = json.Unmarshal(file, &config)

	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	Token = config.Token
	BotPrefix = config.BotPrefix

	return nil

}

var BotId string
var goBot *discordgo.Session

func Start() {
	goBot, err := discordgo.New("Bot " + config.Token)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	u, err := goBot.User("@me")

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	BotId = u.ID

	goBot.AddHandler(messageHandler)

	err = goBot.Open()

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("Bot is running !")
}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == BotId {
		return
	}

	if m.Content == BotPrefix+"ping" {
		_, _ = s.ChannelMessageSend(m.ChannelID, "pong")
	}

	if strings.HasPrefix(m.Content, "!lyrics") {
		fmt.Printf("Command lyrics berjalan!")
		splitMessage := strings.Split(m.Content, " ")
		l := lyrics.New(lyrics.WithAllProviders(), lyrics.WithGeniusLyrics("iCwVYrcAiRXaquxA8gs7F_koVqQcWRjhBsBT5wxli1Pw8jWRiMBXzHbucYpFoZqM"))
		artist := ""
		song := ""
		// mencari index judul lagu dan artist
		indexSong := 0
		indexArtist := 0
		for i := 0; i < len(splitMessage); i++ {
			if splitMessage[i] == "artist" {
				indexArtist = i
			}
			if splitMessage[i] == "song" {
				indexSong = i
			}
		}
		// loop artist
		for i := indexArtist + 1; i < indexSong; i++ {
			artist = artist + " " + splitMessage[i]
		}
		// loop song
		for i := indexSong + 1; i < len(splitMessage); i++ {
			song = song + " " + splitMessage[i]
		}

		lyric, err := l.Search(artist, song)

		if err != nil {
			res := fmt.Sprintf("Lyrics for %v -%v were not found", artist, song)
			_, _ = s.ChannelMessageSend(m.ChannelID, res)
			fmt.Printf("Lyrics for %v -%v were not found", artist, song)
		}
		_, _ = s.ChannelMessageSend(m.ChannelID, lyric)
	}
}

func main() {
	err := ReadConfig()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	Start()

	<-make(chan struct{})
	return
}
