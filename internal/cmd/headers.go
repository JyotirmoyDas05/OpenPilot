package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/JyotirmoyDas05/openpilot/internal/config"
	"github.com/JyotirmoyDas05/openpilot/internal/llm/tools"
	"github.com/spf13/cobra"
)

var headersCmd = &cobra.Command{
	Use:   "headers",
	Short: "Print active HTTP header set",
	Long:  "Displays the User-Agent and resolved provider headers (excluding secrets).",
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, _ := cmd.Flags().GetString("cwd")
		if cwd == "" {
			var err error
			cwd, err = os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to determine cwd: %w", err)
			}
		}
		cfg, err := config.Init(cwd, false)
		if err != nil {
			return fmt.Errorf("failed to init configuration: %w", err)
		}

		ua := tools.UserAgent()
		output := map[string]any{
			"user_agent": ua,
			"providers":  map[string]map[string]string{},
		}
		providers := output["providers"].(map[string]map[string]string)
		for id, p := range cfg.Providers.Seq2() {
			headers := map[string]string{}
			for k, v := range p.ExtraHeaders {
				lk := strings.ToLower(k)
				if lk == "authorization" || lk == "x-api-key" || lk == "api-key" {
					headers[k] = "[REDACTED]"
					continue
				}
				headers[k] = v
			}
			providers[id] = headers
		}
		// Stable ordering for reproducibility
		buf, _ := json.MarshalIndent(output, "", "  ")
		fmt.Fprintln(os.Stdout, string(buf))
		return nil
	},
}
