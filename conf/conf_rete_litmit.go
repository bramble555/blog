package config

type RateLimit struct {
	Enable                 bool    `mapstructure:"enable"`
	Store                  string  `mapstructure:"store"`
	GlobalRPS              float64 `mapstructure:"global_rps"`
	GlobalBurst            int     `mapstructure:"global_burst"`
	IPRPS                  float64 `mapstructure:"ip_rps"`
	IPBurst                int     `mapstructure:"ip_burst"`
	IPTTLSeconds           int     `mapstructure:"ip_ttl_seconds"`
	CleanupIntervalSeconds int     `mapstructure:"cleanup_interval_seconds"`
}
