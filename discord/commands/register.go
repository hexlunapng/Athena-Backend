package commands

import (
	"Athena-Backend/database/models"
	"Athena-Backend/src/profile"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

var RegisterCommand = &discordgo.ApplicationCommand{
	Name:        "register",
	Description: "Register a new Athena account",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "username",
			Description: "Your desired username",
			Required:    true,
		},
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "email",
			Description: "Your email address",
			Required:    true,
		},
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "password",
			Description: "Your password",
			Required:    true,
		},
	},
}

func RegisterCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	logger, _ := zap.NewProduction()

	options := i.ApplicationCommandData().Options
	username := options[0].StringValue()
	email := options[1].StringValue()
	password := options[2].StringValue()
	discordID := i.Member.User.ID

	exists, err := models.UserExists(email)
	if err != nil {
		respondEmbed(s, i, "❌ Error checking existing user: "+err.Error(), true)
		return
	}
	if exists {
		respondEmbed(s, i, "❌ Account with that email already exists.", true)
		return
	}

	accountID := uuid.New().String()

	user := models.UserAccount(accountID, username, email, password, &discordID)
	err = user.Save()
	if err != nil {
		respondEmbed(s, i, "❌ Failed to save user: "+err.Error(), true)
		return
	}

	_, err = profile.CreateProfile(accountID, username, logger)
	if err != nil {
		respondEmbed(s, i, "❌ Failed to create profiles: "+err.Error(), true)
		return
	}

	embed := &discordgo.MessageEmbed{
		Title:       "Account Registered!",
		Description: "**Username:** " + username + "\n**AccountID:** " + accountID,
		Color:       0x00ff00,
	}
	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})
}

func respondEmbed(s *discordgo.Session, i *discordgo.InteractionCreate, msg string, ephemeral bool) {
	flags := discordgo.MessageFlags(0)
	if ephemeral {
		flags = discordgo.MessageFlagsEphemeral
	}
	embed := &discordgo.MessageEmbed{
		Description: msg,
		Color:       0xff0000,
	}
	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
			Flags:  flags,
		},
	})
}
