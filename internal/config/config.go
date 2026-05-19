package config

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

// embedded is set at startup by main.go (via SetEmbedded) so the package
// can stay decoupled from where the config file lives in the source tree.
var embedded []byte

func SetEmbedded(b []byte) { embedded = b }

type Branding struct {
	LauncherName string `toml:"launcher_name"`
	WindowTitle  string `toml:"window_title"`
	PrimaryColor string `toml:"primary_color"`
	LogoPath     string `toml:"logo_path"`
}

type Server struct {
	ID          string `toml:"id"           json:"id"`
	Name        string `toml:"name"         json:"name"`
	ManifestURL string `toml:"manifest_url" json:"manifestUrl"`
	Website     string `toml:"website"      json:"website"`
	NewsFeedURL string `toml:"news_feed_url" json:"newsFeedUrl"`
}

type Security struct {
	ManifestPubkeyHex string `toml:"manifest_pubkey_hex"`
}

type Paths struct {
	ProfilesSubdir string `toml:"profiles_subdir"`
}

type Config struct {
	Branding Branding `toml:"branding"`
	Servers  []Server `toml:"servers"`
	Security Security `toml:"security"`
	Paths    Paths    `toml:"paths"`
}

func Load() (*Config, error) {
	if len(embedded) == 0 {
		return nil, fmt.Errorf("config not embedded — call config.SetEmbedded() in main before Load")
	}
	var c Config
	if err := toml.Unmarshal(embedded, &c); err != nil {
		return nil, fmt.Errorf("parse embedded config: %w", err)
	}
	if len(c.Servers) == 0 {
		return nil, fmt.Errorf("config.toml: at least one [[servers]] entry required")
	}
	return &c, nil
}
