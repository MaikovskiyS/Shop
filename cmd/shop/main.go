package main

import (
	"log"
	"myproject/internal/app"
	"myproject/internal/config"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}
	err = app.RunMigrations(cfg.Psql)
	if err != nil {
		log.Fatal(err)
	}

	err = app.Run(cfg)
	if err != nil {
		log.Fatal(err)
	}
}

/*
logger
linters
gRPC
graphQl
mongo +
metrics +
swagger
ci/cd
werf
async +
clickhouse
elastic
*/
