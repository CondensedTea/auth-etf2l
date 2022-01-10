package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"ratiod/pkg/store"

	"github.com/gofiber/fiber/v2"
)

const discordApiUrl = "https://discord.com/api/users/@me"

func (h *Handler) AuthDiscord(c *fiber.Ctx) error {
	state := c.Query("state")
	if state == "" {
		return fiber.NewError(http.StatusBadRequest, "State is not provided")
	}
	return c.Redirect(h.discordConfig.AuthCodeURL(state), http.StatusFound)
}

func (h *Handler) DiscordCallback(c *fiber.Ctx) error {
	state := c.FormValue("state")

	_, err := h.store.Get(state)
	if err == store.ErrNotFound {
		return fiber.NewError(http.StatusBadRequest, "failed to validate state")
	}
	token, err := h.discordConfig.Exchange(c.Context(), c.FormValue("code"))
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, fmt.Sprintf("failed to get auth token: %v", err))
	}
	res, err := h.discordConfig.Client(c.Context(), token).Get(discordApiUrl)
	if err != nil || res.StatusCode != http.StatusOK {
		return fiber.NewError(res.StatusCode, fmt.Sprintf("failed to get discord api response: %v", err))
	}
	defer res.Body.Close()

	var user DiscordUser

	if err = json.NewDecoder(res.Body).Decode(&user); err != nil {
		return fiber.NewError(http.StatusInternalServerError,
			fmt.Sprintf("failed to decode discord api response: %v", err))
	}

	nameWithTag := fmt.Sprintf("%s#%s", user.Username, user.Discriminator)

	if err = h.store.Set(state, nameWithTag, user.Id, store.Discord); err != nil {
		return fiber.NewError(http.StatusInternalServerError, "failed to save discord name to store")
	}
	return c.Redirect("/")
}
