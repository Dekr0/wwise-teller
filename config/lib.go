package config

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"slices"

	"github.com/Dekr0/wwise-teller/utils"
)

const path = "config.json"

type Config struct {
	DefaultSave string `json:"defaultSave"`
	HelldiversData string `json:"helldiversData"`
	Home string `json:"home"`
	Bookmark []string `json:"bookmark"`
}

func initHome() (string, error) {
	home, err := utils.GetHome()
	if err != nil {
		home, err := os.Executable()
		if err != nil {
			return "", err
		}
		home = filepath.Dir(home)
	}
	return home, nil
}

func New() (*Config, error) {
	home, err := initHome()
	if err != nil {
		return nil, err
	}
	return &Config{
		HelldiversData: home, Home: home, DefaultSave: home, Bookmark: []string{},
	}, nil
}

func Scratch() (*Config, error) {
	c, err := New()
	if err != nil {
		return nil, err
	}
	return c, c.Save()
}

func Load() (*Config, error) {
	_, err := os.Lstat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return Scratch()
		}
		return nil, err
	}
	
	reader, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	var c Config
	if err = json.NewDecoder(reader).Decode(&c); err != nil {
		return nil, err
	}
	return &c, c.Check()
}

func (c *Config) Save() error {
	var data []byte

	data, err := json.MarshalIndent(c, "", "    ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0666)
}

func (c *Config) Check() error {
	_, err := os.Lstat(c.Home)
	if err != nil {
		slog.Error(
			fmt.Sprintf("Invalid home directory: %s", c.Home),
			"error", err,
		)
		slog.Warn("Attempt to fix home directory in config.json")
		c.Home, err = initHome()
		if err != nil {
			return err
		}
	}
	_, err = os.Lstat(c.DefaultSave)
	if err != nil {
		c.DefaultSave = c.Home
		slog.Error(fmt.Sprintf(
			"Invalid default directory for save file dialog : %s", c.Home),
			"error", err,
		)
		slog.Warn(fmt.Sprintf("Setting default directory for save file dialog to %s", c.Home))
	}
	_, err = os.Lstat(c.HelldiversData)
	if err != nil {
		slog.Error(fmt.Sprintf(
				"Invalid Helldivers 2 data directory: %s", c.HelldiversData,
			),
			"error", err,
		)
		slog.Warn(fmt.Sprintf("Setting Helldivers 2 data directory to %s", c.Home))
		c.HelldiversData = c.Home
	}

	clean := []string{}
	for _, s := range c.Bookmark {
		_, err = os.Lstat(s)
		if err != nil {
			continue
		}
		if slices.Contains(clean, s) {
			continue
		}
		clean = append(clean, s)
	}
	c.Bookmark = clean

	return nil
}

func (c *Config) SetHome(home string) error {
	_, err := os.Lstat(home)
	if err == nil {
		c.Home = home
	}
	return err
}

func (c *Config) SetHelldiversData(data string) error {
	_, err := os.Lstat(data)
	if err == nil {
		c.HelldiversData = data 
	}
	return err
}
