package config

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"slices"

	"github.com/Dekr0/wwise-teller/db"
	"github.com/Dekr0/wwise-teller/utils"
)

const path = "config.json"

type Config struct {
	DefaultSave    string `json:"defaultSave"`
	HelldiversData string `json:"helldiversData"`
	IdDatabase     string `json:"idDatabase"`
	Home           string `json:"home"`
	Bookmark     []string `json:"bookmark"`
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
	stat, err := os.Lstat(c.Home)
	if err != nil {
		slog.Error(
			fmt.Sprintf("%s is not a valid directory path.", c.Home),
			"error", err,
		)
		slog.Warn("Attempting to fix home directory in config.json...")
		c.Home, err = initHome()
		if err != nil {
			return err
		}
	}
	if !stat.IsDir() {
		slog.Error(fmt.Sprintf("%s is not a directory.", c.Home))
		slog.Warn("Attempting to fix home directory in config.json...")
		c.Home, err = initHome()
		if err != nil {
			return err
		}
	}

	stat, err = os.Lstat(c.DefaultSave)
	if err != nil {
		slog.Error(fmt.Sprintf(
			"%s is not a valid default directory path for save file dialog.", c.DefaultSave),
			"error", err,
		)
		c.DefaultSave = c.Home
		slog.Warn(fmt.Sprintf("Set default directory for save file dialog to %s", c.Home))
	}
	if !stat.IsDir() {
		slog.Error(fmt.Sprintf("Default path for save file dialog %s is not a directory.", c.DefaultSave))
		c.DefaultSave = c.Home
		slog.Warn(fmt.Sprintf("Set default directory for save file dialog to %s", c.Home))
	}

	stat, err = os.Lstat(c.HelldiversData)
	if err != nil {
		slog.Error(fmt.Sprintf("%s is not a valid directory path for Helldivers 2 data directory.", c.HelldiversData), "error", err)
		c.HelldiversData = c.Home
		slog.Warn(fmt.Sprintf("Set Helldivers 2 data directory to %s", c.Home))
	}
	if !stat.IsDir() {
		slog.Error(fmt.Sprintf("%s is not a directory.", c.HelldiversData))
		c.DefaultSave = c.Home
		slog.Warn(fmt.Sprintf("Set Helldivers 2 data directory to %s", c.Home))
	}

	stat, err = os.Lstat(c.IdDatabase)
	if err != nil {
		slog.Error(fmt.Sprintf("%s is not valid file path for ID database.", c.IdDatabase))
		c.IdDatabase = ""
		slog.Warn("Some specific operations will error since ID database is missing.")
	}
	if !stat.IsDir() {
		slog.Error(fmt.Sprintf("%s is not a file.", c.IdDatabase))
		c.IdDatabase = ""
		slog.Warn("Some specific operations will error since ID database is missing.")
	}

	if err := os.Setenv(db.DatabaseEnv, c.IdDatabase); err != nil {
		slog.Error(fmt.Sprintf("Failed to set %s enviromental variable using %s", db.DatabaseEnv, c.IdDatabase))
		c.IdDatabase = ""
		slog.Warn("Some specific operations will error since ID database is missing.")
	}

	clean := make([]string, 0, len(c.Bookmark))
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
