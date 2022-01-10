package store

import "fmt"

var (
	ErrNotFound    = fmt.Errorf("state is not found in storage")
	ErrIsCompleted = fmt.Errorf("this player already authorized in system")
)

type Provider int

const (
	Steam Provider = iota
	Discord
)

type Platforms struct {
	SteamName   string `json:"steam_name"`
	SteamID     string `json:"steam_id"`
	DiscordName string `json:"discord_name"`
	DiscordID   string `json:"discord_id"`
	completed   bool
}

type InMemoryStore struct {
	_map map[string]*Platforms
}

func NewStore() *InMemoryStore {
	return &InMemoryStore{
		_map: make(map[string]*Platforms),
	}
}

func (s *InMemoryStore) Init(key string) {
	s._map[key] = &Platforms{}
}

func (s *InMemoryStore) Set(key, name, id string, prov Provider) error {
	_, ok := s._map[key]
	if !ok {
		return ErrNotFound
	}
	if s._map[key].completed == true {
		return ErrIsCompleted
	}
	switch prov {
	case Steam:
		s._map[key].SteamName = name
		s._map[key].SteamID = id
	case Discord:
		s._map[key].DiscordName = name
		s._map[key].DiscordID = id
	}
	s.updateIfCompleted(key)
	return nil
}

func (s InMemoryStore) Get(key string) (*Platforms, error) {
	names, ok := s._map[key]
	if !ok {
		return nil, ErrNotFound
	}
	if names.completed {
		return names, ErrIsCompleted
	}
	return s._map[key], nil
}

func (s *InMemoryStore) updateIfCompleted(key string) {
	_, ok := s._map[key]
	if !ok {
		return
	}
	if s._map[key].SteamName != "" && s._map[key].DiscordName != "" {
		s._map[key].completed = true
	}
}
