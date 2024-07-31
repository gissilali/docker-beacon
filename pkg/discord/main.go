package discord

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"os"
)

type DockerBot struct {
	session *discordgo.Session
}

func OpenSession() (*DockerBot, error) {
	token := os.Getenv("DISCORD_BOT_TOKEN")
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Printf("error creating docker session: %v\n", err)
		return nil, err
	}

	session.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsDirectMessages + discordgo.IntentGuildMessages)
	err = session.Open()

	if err != nil {
		fmt.Printf("error opening connection: %v\n", err)
		return nil, err
	}

	fmt.Println("Connected to Discord server")

	return &DockerBot{
		session: session,
	}, nil
}

func (bot DockerBot) SendMessage(message string) (*discordgo.Message, error) {
	send, err := bot.session.ChannelMessageSend("1265637759557832778", message)
	return send, err
}
