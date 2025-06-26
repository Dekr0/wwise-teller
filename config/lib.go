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

const DefaultConfigPath = "config.json"

type Config struct {
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
		Home: home, Bookmark: []string{},
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
	_, err := os.Lstat(DefaultConfigPath)
	if err != nil {
		if os.IsNotExist(err) {
			return Scratch()
		}
		return nil, err
	}
	
	reader, err := os.Open(DefaultConfigPath)
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

	return os.WriteFile(DefaultConfigPath, data, 0666)
}

func (c *Config) Check() error {
	err := c.CheckHome()
	if err != nil {
		return err
	}
	c.CheckIdDatabase()

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

func (c *Config) CheckHome() error {
	if !filepath.IsAbs(c.Home) {
		slog.Error(fmt.Sprintf("%s is not an absolute path", c.Home))
		slog.Warn("Attempting to fix home directory in config.json...")

		var err error
		c.Home, err = initHome()
		if err != nil {
			return err
		}
		return nil
	}

	stat, err := os.Lstat(c.Home)
	if err != nil {
		slog.Error(fmt.Sprintf("%s is not a valid directory path.", c.Home), "error", err)
		slog.Warn("Attempting to fix home directory in config.json...")
		c.Home, err = initHome()
		if err != nil {
			return err
		}
		return nil
	}

	if !stat.IsDir() {
		slog.Error(fmt.Sprintf("%s is not a directory.", c.Home))
		slog.Warn("Attempting to fix home directory in config.json...")
		c.Home, err = initHome()
		if err != nil {
			return err
		}
		return nil
	}

	return nil
}

func (c *Config) SetHome(home string) error {
	if !filepath.IsAbs(home) {
		return fmt.Errorf("%s is not an absolute directory", home)
	}
	_, err := os.Lstat(home)
	if err == nil {
		c.Home = home
	}
	return err
}

func (c *Config) CheckIdDatabase() {
	if !filepath.IsAbs(c.IdDatabase) {
		slog.Error(fmt.Sprintf("%s is not an absolute path.", c.IdDatabase))
		c.IdDatabase = ""
		slog.Warn("Some specific operations will error since ID database is missing.")
		return
	}
	stat, err := os.Lstat(c.IdDatabase)
	if err != nil {
		slog.Error(fmt.Sprintf("%s is not a valid file path for ID database.", c.IdDatabase))
		c.IdDatabase = ""
		slog.Warn("Some specific operations will error since ID database is missing.")
		return
	}
	if !stat.IsDir() {
		slog.Error(fmt.Sprintf("%s is not a file.", c.IdDatabase))
		c.IdDatabase = ""
		slog.Warn("Some specific operations will error since ID database is missing.")
		return
	}
	if err := os.Setenv(db.DatabaseEnv, c.IdDatabase); err != nil {
		slog.Error(fmt.Sprintf("Failed to set %s enviromental variable using %s", db.DatabaseEnv, c.IdDatabase))
		c.IdDatabase = ""
		slog.Warn("Some specific operations will error since ID database is missing.")
		return
	}
	if err := db.CheckDatabaseEnv(); err != nil {
		slog.Error(err.Error())
		slog.Warn("Some specific operations will error since ID database is missing.")
		return
	}
}
