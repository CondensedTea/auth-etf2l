package main

import (
	"log"
	"math/rand"
	"auth-etf2l/pkg/config"
	"auth-etf2l/pkg/discord"
	"auth-etf2l/pkg/handler"
	"auth-etf2l/pkg/store"
	"time"

	"github.com/gofiber/fiber/v2/middleware/logger"

	"github.com/gofiber/fiber/v2"
)

var seed = time.Now().Unix()

func main() {
	rand.Seed(seed)

	app := fiber.New()
	app.Use(logger.New())

	cfg := config.NewEnvConfig()
	s := store.NewStore()
	d, err := discord.NewDiscord(cfg.DiscordBotToken(), cfg.RoleID(), cfg.GuildId())
	if err != nil {
		log.Fatalf("failed to init discord bot: %v", err)
	}
	h := handler.NewHandler(cfg, s, d.UpdatesChannel())

	app.Static("/css", "./src/css")
	app.Static("/js", "./src/js")
	app.Static("/", "./src/html")

	app.Get("/discord/callback", h.DiscordCallback) // todo consistent urls
	app.Get("/auth/discord", h.AuthDiscord)

	app.Get("/steam/callback", h.SteamCallback)
	app.Get("/auth/steam/", h.AuthSteam)

	app.Post("/api/register", h.Register)
	app.Get("/api/:state", h.CheckState)

	if err = app.Listen(":3000"); err != nil {
		log.Fatal(err)
	}
}
