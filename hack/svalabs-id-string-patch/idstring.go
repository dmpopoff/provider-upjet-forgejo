package provider

import "strconv"

func mustParseInt64(s string) int64 {
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	}
	return v
}

// Empty or "0" means upjet has not yet received a provider-allocated id
// (Observe before Create). Treat as resource absent.
func isUnsetIDString(s string) bool {
	return s == "" || s == "0"
}
