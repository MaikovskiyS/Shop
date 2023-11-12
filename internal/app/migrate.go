package app

import (
	"log"
	"myproject/internal/config"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(cfg *config.Postgres) error {
	dr, err := migrate.New(cfg.MigrationPath, cfg.ConnString())
	if err != nil {
		return err
	}
	defer dr.Close()

	if err := dr.Up(); err != nil && err != migrate.ErrNoChange {
		log.Println(err)
		dr.Drop()
		return err
	}

	//dr.Drop()
	log.Println("migrations done")
	return nil
}
