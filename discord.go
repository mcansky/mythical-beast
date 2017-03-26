package main

import (
    "fmt"
    "github.com/bwmarrin/discordgo"
)
// Discord messenger

func DiscordMessage(channel_name string, message string) {
	discord, err := discordgo.New("Bot " + Token)
	if err != nil {
		console_logger("ERROR", "discordMessage", fmt.Sprintf("error creating Discord session %s", err))
		return
	}

	// Get the account information.
	u, err := discord.User("@me")
	if err != nil {
		fmt.Println("error obtaining account details,", err)
	}

	// Store the account ID for later use.
	BotID = u.ID

	_, err = discord.ChannelMessageSend(ChannelID, message)
	if err != nil {
		fmt.Println("could not send message to discord,", err)
	}

	return
}

func register_bot() {
	discord, err := discordgo.New("Bot " + Token)
	if err != nil {
		console_logger("ERROR", "discordMessage", fmt.Sprintf("error creating Discord session %s", err))
		return
	}

	// Get the account information.
	u, err := discord.User("@me")
	if err != nil {
		fmt.Println("error obtaining account details,", err)
	}

	// Store the account ID for later use.
	BotID = u.ID

	// Register basics as a callback for the messageCreate events.
	discord.AddHandler(basics)

	// Open the websocket and begin listening.
	err = discord.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}
}
