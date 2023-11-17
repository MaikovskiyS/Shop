package app

import (
	"log"
	"myproject/internal/config"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// Run migrations
func RunMigrations(cfg *config.Postgres) error {
	dr, err := migrate.New(cfg.MigrationPath, cfg.ConnString())
	if err != nil {
		return err
	}
	defer dr.Close()

	if err := dr.Up(); err != nil && err != migrate.ErrNoChange {
		log.Println(err)
		err = dr.Drop()
		if err != nil {
			return err
		}
		return err
	}

	//dr.Drop()
	log.Println("migrations done")
	return nil
}
