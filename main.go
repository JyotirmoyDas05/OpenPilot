package main

import (
	"log/slog"
	"net/http"
	"os"

	_ "net/http/pprof" // profiling

	_ "github.com/joho/godotenv/autoload" // automatically load .env files

	"github.com/JyotirmoyDas05/openpilot/internal/cmd"
	"github.com/JyotirmoyDas05/openpilot/internal/log"
)

func main() {
	defer log.RecoverPanic("main", func() {
		slog.Error("Application terminated due to unhandled panic")
	})

	// Enable pprof when OPENPILOT_PROFILE is set.
	profileEnv := os.Getenv("OPENPILOT_PROFILE")
	if profileEnv != "" {
		go func() {
			slog.Info("Serving pprof at localhost:6060")
			if httpErr := http.ListenAndServe("localhost:6060", nil); httpErr != nil {
				slog.Error("Failed to pprof listen", "error", httpErr)
			}
		}()
	}

	cmd.Execute()
}
