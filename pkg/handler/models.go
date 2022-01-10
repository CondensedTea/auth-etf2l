package handler

type RegisterResponse struct {
	State string `json:"state"`
}

type DiscordUser struct {
	Id            string      `json:"id"`
	Username      string      `json:"username"`
	Avatar        string      `json:"avatar"`
	Discriminator string      `json:"discriminator"`
	PublicFlags   int         `json:"public_flags"`
	Flags         int         `json:"flags"`
	Banner        interface{} `json:"banner"`
	BannerColor   string      `json:"banner_color"`
	AccentColor   int         `json:"accent_color"`
	Locale        string      `json:"locale"`
	MfaEnabled    bool        `json:"mfa_enabled"`
	PremiumType   int         `json:"premium_type"`
}

type Etf2lUser struct {
	Player struct {
		Bans       interface{} `json:"bans"`
		Classes    []string    `json:"classes"`
		Country    string      `json:"country"`
		Id         int         `json:"id"`
		Name       string      `json:"name"`
		Registered int         `json:"registered"`
		Steam      struct {
			Avatar string `json:"avatar"`
			Id     string `json:"id"`
			Id3    string `json:"id3"`
			Id64   string `json:"id64"`
		} `json:"steam"`
	} `json:"player"`
}
