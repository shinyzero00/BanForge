package config

type Firewall struct {
	Name    string `toml:"name"`
	Config  string `toml:"config"`
	BanTime int    `toml:"ban_time"`
}

type Service struct {
	Name    string `toml:"name"`
	LogPath string `toml:"log_path"`
	Enabled bool   `toml:"enabled"`
}

type Config struct {
	Firewall Firewall `toml:"firewall"`
	Service  Service  `toml:"service"`
}
