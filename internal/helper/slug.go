package helper

import "strings"

func GenerateSlug(name string) string {
	name = strings.ToLower(strings.TrimSpace(name))

	var builder strings.Builder
	lastDash := false

	for _, r := range name {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			builder.WriteRune(r)
			lastDash = false
			continue
		}

		if !lastDash && builder.Len() > 0 {
			builder.WriteRune('-')
			lastDash = true
		}
	}

	return strings.TrimSuffix(builder.String(), "-")
}
