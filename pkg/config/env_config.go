package config

import "os"

type EnvConfig struct {
	baseUrl             string
	discordClientID     string
	discordClientSecret string
	discordBotToken     string
	unauthorizedRoleID  string
	discordGuildId      string
}

func NewEnvConfig() *EnvConfig {
	return &EnvConfig{
		baseUrl:             os.Getenv("BASE_URL"),
		discordClientID:     os.Getenv("CLIENT_ID"),
		discordClientSecret: os.Getenv("CLIENT_SECRET"),
		discordBotToken:     os.Getenv("BOT_TOKEN"),
		unauthorizedRoleID:  os.Getenv("ROLE_ID"),
		discordGuildId:      os.Getenv("GUILD_ID"),
	}
}

func (c EnvConfig) URL() string {
	return c.baseUrl
}

func (c EnvConfig) DiscordClientID() string {
	return c.discordClientID
}

func (c EnvConfig) DiscordClientSecret() string {
	return c.discordClientSecret
}

func (c EnvConfig) DiscordBotToken() string {
	return c.discordBotToken
}

func (c EnvConfig) RoleID() string {
	return c.unauthorizedRoleID
}

func (c EnvConfig) GuildId() string {
	return c.discordGuildId
}
