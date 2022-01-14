package handler

import (
	"fmt"
	"math/rand"
	"net/http"
	"auth-etf2l/pkg/store"

	"github.com/gofiber/fiber/v2"
	"github.com/ravener/discord-oauth2"
	"golang.org/x/oauth2"
)

type Configurator interface {
	URL() string
	DiscordClientID() string
	DiscordClientSecret() string
}

type Store interface {
	Init(key string)
	Set(key, name, id string, prov store.Provider) error
	Get(key string) (*store.Platforms, error)
}

type Handler struct {
	store          Store
	discordConfig  oauth2.Config
	baseURL        string
	updatesChannel chan string
}

func NewHandler(cfg Configurator, store Store, ch chan string) *Handler {
	return &Handler{
		store: store,
		discordConfig: oauth2.Config{
			ClientID:     cfg.DiscordClientID(),     //"857714670570962944", serviceURL = "http://localhost:3000"
			ClientSecret: cfg.DiscordClientSecret(), // "90GLBttHR9-zMhJWDQMMdH4gj7RxC3nk",
			Scopes:       []string{discord.ScopeIdentify},
			Endpoint:     discord.Endpoint,
		},
		baseURL:        cfg.URL(),
		updatesChannel: ch,
	}
}

func (h *Handler) Register(c *fiber.Ctx) error {
	const length = 64
	b := make([]byte, length)
	rand.Read(b)
	state := fmt.Sprintf("%x", b)[:length]

	h.store.Init(state)

	return c.JSON(RegisterResponse{
		state,
	})
}

func (h *Handler) CheckState(c *fiber.Ctx) error {
	state := c.Params("state")
	names, err := h.store.Get(state)
	if err == store.ErrNotFound {
		return fiber.NewError(http.StatusNotFound, "unknown state")
	}
	if err == store.ErrIsCompleted {
		h.updatesChannel <- names.DiscordID
	}
	return c.JSON(names)
}
