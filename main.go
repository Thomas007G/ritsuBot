package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/KurozeroPB/kitsu-go"
	"github.com/bwmarrin/discordgo"
)

// Variables used for command line parameters
var (
	Token   = "insert token here"
	dg, err = discordgo.New("Bot " + Token)
)

func init() {

	flag.StringVar(&Token, "t", Token, "Bot Token")
	flag.Parse()
}

func main() {

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)
	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	dg.UpdateStatus(1, "Type -help")
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	fmt.Println("------------------------------------------")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}
	// Says hi to the user
	if strings.HasPrefix(m.Content, "-") {
		// if strings.Contains(m.Content, "help") && len(m.Mentions) == 1 {
		//^alllows mention of users
		if strings.Contains(m.Content, "help") {
			s.ChannelMessageSend(m.ChannelID, "Type `-` followed by:")
			s.ChannelMessageSend(m.ChannelID, "- `-kitsu` `anime` --> Searches Kitsu database for relevent shows.")
			s.ChannelMessageSend(m.ChannelID, "- `-kitsu` `manga` --> Searches Kitsu database for relevent manga.")
			s.ChannelMessageSend(m.ChannelID, "- `-kitsu` `user` --> Returns information about Kitsu user.")
			s.ChannelMessageSend(m.ChannelID, "- `-kitsu` `char` --> Returns information about relevant character.")
		}

		if strings.Contains(m.Content, "hi") {
			s.ChannelMessageSend(m.ChannelID, "**"+"Hello "+m.Author.Mention()+" :wave:**")
		}
		// If the message is "ping" reply with "Pong!"
		if strings.Contains(m.Content, "ping") {
			s.ChannelMessageSend(m.ChannelID, "Pong!")
		}

		// If the message is "pong" reply with "Ping!"
		if strings.Contains(m.Content, "pong") {
			s.ChannelMessageSend(m.ChannelID, "Ping!")
		}
		if strings.Contains(m.Content, "kitsu anime") {
			var query = strings.Split(m.Content, "anime ")[1]
			fmt.Println("searching anime:", query)
			anime, e := kitsu.SearchAnime(query, 0)
			if e != nil {
			}
			if anime.Attributes.PosterImage.Original != "" {
				s.ChannelMessageSend(m.ChannelID, anime.Attributes.PosterImage.Original)
			}
			if anime.Attributes.CanonicalTitle != "" {
				s.ChannelMessageSend(m.ChannelID, "		** "+anime.Attributes.AbbreviatedTitles[0]+"**")
			} else {
				s.ChannelMessageSend(m.ChannelID, "		** "+anime.Attributes.CanonicalTitle+"**")
			}
			if anime.Attributes.Synopsis != "" {
				s.ChannelMessageSend(m.ChannelID, anime.Attributes.Synopsis)
			}
			if anime.Attributes.PopularityRank != 0 {
				s.ChannelMessageSend(m.ChannelID, "Is ranked #"+strconv.Itoa(anime.Attributes.PopularityRank)+" in popularity.")
			}
		}
		if strings.Contains(m.Content, "kitsu manga") {
			var query = strings.Split(m.Content, "manga ")[1]
			//words := strings.Fields(query)
			//for i, word := range words {
			//fmt.Println(i, " => ", word)
			//query = strings.Split(word, " ")
			//query = append(query, word)

			//var query = []string{m.Content}
			fmt.Println("searching manga:", query)
			manga, e := kitsu.SearchManga(query, 0)
			if e != nil {
			}
			s.ChannelMessageSend(m.ChannelID, manga.Attributes.AbbreviatedTitles[0])
			s.ChannelMessageSend(m.ChannelID, manga.Attributes.PosterImage.Original)
			s.ChannelMessageSend(m.ChannelID, manga.Attributes.Synopsis)
			s.ChannelMessageSend(m.ChannelID, "Is ranked #"+strconv.Itoa(manga.Attributes.PopularityRank)+" in popularity.")
		}
		if strings.Contains(m.Content, "kitsu user") {
			//var query = []string{m.Content}
			var query = strings.Split(m.Content, "user ")[1]
			fmt.Println("searching user: ", query)
			user, e := kitsu.SearchUser(query)
			if e != nil {
			}
			//var embed = user.Attributes.Avatar.Small
			//var name *discordgo.MessageEmbed
			//name.URL = embed
			//s.ChannelEmbed(name.URL)

			//user.ID = int(user.ID)
			s.ChannelMessageSend(m.ChannelID, string(user.Attributes.Avatar.Small))
			s.ChannelMessageSend(m.ChannelID, ("User: " + user.Attributes.Name + " has a" + user.Attributes.WaifuOrHusbando +
				" and has watched " + strconv.Itoa(user.Attributes.LifeSpentOnAnime) + " hours of anime"))
		}
		//s.ChannelMessageSend(m.ChannelID, string(kitsu.GetStats(ID)))
		if strings.Contains(m.Content, "kitsu char") {
			var query = strings.Split(m.Content, "char ")[1]
			fmt.Println("searching char: ", query)
			char, e := kitsu.SearchCharacter(query)
			s.ChannelMessageSend(m.ChannelID, char.Attributes.Image.Original)
			s.ChannelMessageSend(m.ChannelID, char.Attributes.Name)
			s.ChannelMessageSend(m.ChannelID, char.Attributes.Description)
			if e != nil {
			}
		}
	}
}
