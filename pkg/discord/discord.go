package discord

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

type Discord struct {
	bot     *discordgo.Session
	updates chan string
	guildId string
	roleId  string
}

var roleID string

func NewDiscord(token, roleId, guildId string) (*Discord, error) {
	roleID = roleId

	discord, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}
	fmt.Println("discord is running !")

	discord.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMembers

	discord.AddHandler(ready)
	discord.AddHandler(addRoleOnServerJoin)

	err = discord.Open()
	if err != nil {
		fmt.Println("Error opening Discord session: ", err)
	}

	updatesChan := make(chan string)

	d := &Discord{
		bot:     discord,
		updates: updatesChan,
		guildId: guildId,
		roleId:  roleID,
	}
	go d.listenForRoleUpdates(updatesChan)

	return d, nil
}

func ready(s *discordgo.Session, event *discordgo.Ready) {
	fmt.Println("hello i am running")
}

func addRoleOnServerJoin(s *discordgo.Session, e *discordgo.GuildMemberAdd) {
	fmt.Println("new server member !")
	if err := s.GuildMemberRoleAdd(e.GuildID, e.User.ID, roleID); err != nil {
		fmt.Println(err)
	}
}

func (d *Discord) listenForRoleUpdates(ch chan string) {
	for user := range ch {
		err := d.bot.GuildMemberRoleRemove(d.guildId, user, d.roleId)
		if err != nil {
			log.Println(err)
		}
	}
}

func (d Discord) UpdatesChannel() chan string {
	return d.updates
}
