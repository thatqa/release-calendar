package cli

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"release-calendar/backend/internal/server"
	"release-calendar/backend/internal/store"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run HTTP server",
	RunE: func(cmd *cobra.Command, args []string) error {
		_ = godotenv.Load()
		cfg := store.ConfigFromEnv()
		db, err := store.Open(cfg)
		if err != nil {
			return err
		}

		r := server.Router(db)
		addr := fmt.Sprintf(":%s", os.Getenv("BACKEND_PORT"))
		if addr == ":" {
			addr = ":8080"
		}
		log.Printf("HTTP server on %s", addr)
		return r.Run(addr)
	},
}
