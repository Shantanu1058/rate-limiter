package limiter

import "fmt"

func BuildGlobalKey(method, path string) string {
	return fmt.Sprintf("leaky_bucket:endpoint:%s:%s", method, path)
}

func BuildIdentityKey(identity, method, path string) string {
	return fmt.Sprintf("leaky_bucket:identity:%s:%s:%s", identity, method, path)
}
