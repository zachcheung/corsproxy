package corsproxy

import "strings"

// wildcard represents a pattern for wildcard matching in allowed targets.
type wildcard struct {
	prefix string
	suffix string
}

// match checks if the target URL matches the wildcard pattern.
func (w wildcard) match(s string) bool {
	if len(s) >= len(w.prefix)+len(w.suffix) {
		if after, found := strings.CutPrefix(s, w.prefix); found {
			return strings.HasSuffix(after, w.suffix) || strings.Contains(after, w.suffix+"/")
		}
	}

	return false
}
