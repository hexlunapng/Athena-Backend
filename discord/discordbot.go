package discord

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const (
	blurple    = "\033[38;5;63m"
	resetColor = "\033[0m"
	discordTag = "[DISCORD]"
)

func colorizeDiscord(text string) string {
	return fmt.Sprintf("%s%s%s", blurple, text, resetColor)
}

func StartAthenaBackendDiscordBot(token string) (*discordgo.Session, error) {
	intents := discordgo.IntentsGuildMessages | discordgo.IntentsDirectMessages | discordgo.IntentsMessageContent

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, fmt.Errorf("%s Error creating Discord session: %w", colorizeDiscord(discordTag), err)
	}

	dg.Identify.Intents = intents
	dg.AddHandler(pingPongHandler)

	err = dg.Open()
	if err != nil {
		return nil, fmt.Errorf("%s Error opening Discord session: %w", colorizeDiscord(discordTag), err)
	}

	fmt.Println(colorizeDiscord(discordTag), "Bot is now running...")
	return dg, nil
}

func pingPongHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	fmt.Printf("%s Message received: %s\n", colorizeDiscord(discordTag), m.Content)

	if strings.EqualFold(m.Content, "!ping") {
		_, err := s.ChannelMessageSend(m.ChannelID, "Pong!")
		if err != nil {
			fmt.Printf("%s Failed to send Pong: %v\n", colorizeDiscord(discordTag), err)
		}
	}
}
