package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"auth-etf2l/pkg/store"
	"regexp"

	"github.com/gofiber/fiber/v2"
)

const (
	steamApiURL = "https://steamcommunity.com/openid/login"
	etf2lApiURL = "https://api.etf2l.org"
)

var steamOpenIdUrlRegexp = regexp.MustCompile(`https://steamcommunity.com/openid/id/(\d+)`)

var ErrUnknownEtf2lUser = fmt.Errorf("unknown ETF2L user")

func (h *Handler) AuthSteam(c *fiber.Ctx) error {
	state := c.Query("state")
	if state == "" {
		return fiber.NewError(http.StatusBadRequest, "State is not provided")
	}
	redirectUrl := constructAuthParams(h.baseURL + "/steam/callback?state=" + state)

	return c.Redirect(redirectUrl, http.StatusFound)
}

func (h *Handler) SteamCallback(c *fiber.Ctx) error {
	state := c.Query("state")
	if state == "" {
		return fiber.NewError(http.StatusBadRequest, "State is not provided")
	}
	claimedIdValue := c.Query("openid.claimed_id")
	claimedIdDecoded, err := url.QueryUnescape(claimedIdValue)
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, "failed to parse query params")
	}
	claimedId := steamOpenIdUrlRegexp.FindStringSubmatch(claimedIdDecoded)
	if len(claimedId) < 1 {
		return fiber.NewError(http.StatusInternalServerError, "failed to extract steamid64 from query params")
	}
	steamId := claimedId[1]

	username, err := checkEtf2lProfile(steamId)
	if err == ErrUnknownEtf2lUser {
		return fiber.NewError(http.StatusNetworkAuthenticationRequired, "etf2l account not found")
	}
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, "failed to retrieve etf2l user")
	}
	err = h.store.Set(state, username, steamId, store.Steam)
	if err == store.ErrNotFound {
		return fiber.NewError(http.StatusNotFound, "unknown state")
	}
	if err == store.ErrIsCompleted {
		return fiber.NewError(http.StatusNotAcceptable, "already authorized on all platforms")
	}
	return c.Redirect("/", http.StatusSeeOther)
}

func constructAuthParams(baseUrl string) string {
	template := "openid.ns=http://specs.openid.net/auth/2.0" +
		"&openid.mode=checkid_setup" +
		"&openid.return_to=%s" +
		"&openid.realm=%s" +
		"&openid.identity=http://specs.openid.net/auth/2.0/identifier_select" +
		"&openid.claimed_id=http://specs.openid.net/auth/2.0/identifier_select"
	params := fmt.Sprintf(template, baseUrl, baseUrl)
	return steamApiURL + "?" + url.PathEscape(params)
}

func checkEtf2lProfile(steamId string) (string, error) {
	etf2lUrl := fmt.Sprintf("%s/player/%s.json", etf2lApiURL, steamId)
	resp, err := http.Get(etf2lUrl)
	if err != nil {
		return "", fmt.Errorf("etf2l api returned status %s, err: %v", resp.Status, err)
	}
	if resp.StatusCode == http.StatusInternalServerError {
		return "", ErrUnknownEtf2lUser
	}
	var user Etf2lUser
	if err = json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return "", fmt.Errorf("failed to decode etf2l api response: %v", err)
	}
	return user.Player.Name, nil
}
