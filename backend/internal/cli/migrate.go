package cli

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"release-calendar/backend/migration"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"
	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run DB migrations",
	RunE: func(cmd *cobra.Command, args []string) error {
		_ = godotenv.Load()
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&multiStatements=true",
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASS"),
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_NAME"),
		)
		db, err := sql.Open("mysql", dsn)
		if err != nil {
			return err
		}
		defer db.Close()

		if err := goose.SetDialect("mysql"); err != nil {
			return err
		}
		goose.SetBaseFS(migration.MigrationsFS)

		if err := goose.Up(db, "."); err != nil {
			return err
		}
		log.Println("Migrations applied successfully")
		return nil
	},
}
