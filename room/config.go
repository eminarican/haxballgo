package room

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

type Config struct {
	Bot struct {
		Active bool
		Name   string
	}
	General struct {
		Name      string
		Token     string
		MaxPlayer int
	}
	Security struct {
		Public   bool
		Password string
	}
}

func (c Config) String() string {
	return `room = top.HBInit({
		maxPlayers: ` + strconv.Itoa(c.General.MaxPlayer) + `,
		noPlayer: ` + strconv.FormatBool(!c.Bot.Active) + `,
		playerName: "` + c.Bot.Name + `",
		public: ` + strconv.FormatBool(c.Security.Public) + `,
		roomName: "` + c.General.Name + `",
		token: "` + c.General.Token + `",
		` + func() string {
		if len(c.Security.Password) == 0 {
			return ``
		} else {
			return `password: "` + c.Security.Password + `"`
		}
	}() + `
	})`
}

func defaultConfig() Config {
	return Config{
		struct {
			Active bool
			Name   string
		}{
			false,
			"Bot",
		},
		struct {
			Name      string
			Token     string
			MaxPlayer int
		}{
			"My Room",
			"",
			16,
		},
		struct {
			Public   bool
			Password string
		}{
			true,
			"",
		},
	}
}

func readConfig() (Config, error) {
	c := defaultConfig()
	if _, err := os.Stat("config.json"); os.IsNotExist(err) {
		data, err := json.MarshalIndent(c, "", "  ")
		if err != nil {
			return c, fmt.Errorf("failed encoding default config: %v", err)
		}
		if err := ioutil.WriteFile("config.json", data, 0644); err != nil {
			return c, fmt.Errorf("failed creating config: %v", err)
		}
		return c, nil
	}
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		return c, fmt.Errorf("error reading config: %v", err)
	}
	if err := json.Unmarshal(data, &c); err != nil {
		return c, fmt.Errorf("error decoding config: %v", err)
	}
	return c, nil
}
