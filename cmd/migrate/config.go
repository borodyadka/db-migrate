package main

import (
	migrate "github.com/borodyadka/db-migrate"
	"github.com/jessevdk/go-flags"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
	"net/url"
	"os"
)

var (
	logLevel     log.Level
	databaseURL  *url.URL
	historyTable *migrate.Table
	sourceDir    string
	command      commandType
)

type specification struct {
	LogLevel     string `envconfig:"LOG_LEVEL" default:"INFO" short:"v" long:"verbosity"`
	DatabaseURL  string `envconfig:"DATABASE_URL" short:"d" long:"db-url"`
	HistoryTable string `envconfig:"HISTORY_TABLE" short:"t" long:"history-table"`
	SourceDir    string `envconfig:"SOURCE_DIR" short:"s" long:"source"`
}

func parseConfig() error {
	_ = godotenv.Load()

	config := new(specification)
	err := envconfig.Process("", config)
	if err != nil {
		return err
	}
	cmd, err := flags.Parse(config)
	if err != nil {
		return err
	}
	command, err = parseCommand(cmd)
	if err != nil {
		return err
	}
	logLevel, err = log.ParseLevel(config.LogLevel)
	if err != nil {
		return err
	}
	log.SetFormatter(&log.TextFormatter{DisableSorting: false})
	log.SetOutput(os.Stdout)
	log.SetLevel(logLevel)

	databaseURL, err = url.Parse(config.DatabaseURL)
	if err != nil {
		return err
	}
	log.Debugf("using driver: %s", databaseURL.Scheme)

	if config.HistoryTable == "" {
		config.HistoryTable = "migrations_history"
	}
	historyTable, err = migrate.ParseTableName(config.HistoryTable, migrate.Dialect(databaseURL.Scheme))
	if err != nil {
		return err
	}
	log.Debugf("using history table: %s", historyTable.String())

	sourceDir = config.SourceDir
	if sourceDir == "" {
		sourceDir, _ = os.Getwd()
	}
	log.Debugf("source directory: %s", sourceDir)

	return nil
}
