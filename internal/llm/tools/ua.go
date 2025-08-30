package tools

import (
	"fmt"

	"github.com/JyotirmoyDas05/openpilot/internal/version"
)

// UserAgent returns the HTTP User-Agent string used for outbound provider/tool requests.
// Centralizing this lets us update branding/version in one place and avoid legacy names.
func UserAgent() string {
	// Format: openpilot/<version>
	return fmt.Sprintf("openpilot/%s", version.Version)
}
