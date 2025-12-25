package policy

import "strings"

type Resolver interface {
	Resolve(method, path string) (LimitTiers, bool)
}

type ExactMatchResolver struct {
	rules map[string]LimitTiers
}

func buildKey(method, path string) string {
	return strings.ToUpper(method) + ":" + path
}

func NewExactMatchResolver(cfg Config) *ExactMatchResolver {
	rules := make(map[string]LimitTiers, len(cfg.Descriptors))

	for _, d := range cfg.Descriptors {
		key := buildKey(d.Match.Method, d.Match.Path)
		rules[key] = d.Limits
	}

	return &ExactMatchResolver{
		rules: rules,
	}
}

func (r *ExactMatchResolver) Resolve(method, path string) (LimitTiers, bool) {
	key := buildKey(method, path)
	limits, ok := r.rules[key]
	return limits, ok
}
