package policy

type RateLimit struct {
	Capacity       int     `yaml:"capacity"`
	LeakRatePerSec float64 `yaml:"leak_rate_per_sec"`
}

type LimitTiers struct {
	Global   *RateLimit `yaml:"global"`
	Identity *RateLimit `yaml:"identity"`
}

type MatchRule struct {
	Method string `yaml:"method"`
	Path   string `yaml:"path"`
}

type Descriptor struct {
	Match  MatchRule  `yaml:"match"`
	Limits LimitTiers `yaml:"limits"`
}

type Config struct {
	Domain      string       `yaml:"domain"`
	Descriptors []Descriptor `yaml:"descriptors"`
}
