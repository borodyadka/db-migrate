package main

import (
	"context"
	"github.com/borodyadka/db-migrate"
	"github.com/borodyadka/db-migrate/drivers/postgres"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
	"time"
)

func init() {
	if err := parseConfig(); err != nil {
		log.Fatal(err)
	}
}

func now() int64 {
	v, _ := strconv.ParseInt(time.Now().Format(migrate.TimestampFormat), 10, 63)
	return v
}

func up(driver migrate.Driver, count int) error {
	fm, err := migrate.ListMigrations(sourceDir)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	dm, err := driver.MigrationsHistory(ctx, historyTable)
	if err != nil {
		return err
	}

	n := 1
	for _, m := range subEntries(fm, dm) {
		if n > count {
			break
		}
		if err := driver.ApplyMigration(context.Background(), historyTable, &m); err != nil {
			return err
		}
		log.WithField("file", m.Filename("up")).Info("applied")
		n++
	}
	return nil
}

func down(driver migrate.Driver, count int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	dm, err := driver.MigrationsHistory(ctx, historyTable)
	if err != nil {
		return err
	}

	ldm := len(dm) - 1
	for i := ldm; i > ldm-count; i-- {
		if i < 0 {
			break
		}
		if err := driver.RevertMigration(context.Background(), historyTable, &dm[i]); err != nil {
			return err
		}
		log.WithField("file", dm[i].Filename("down")).Info("reverted")
	}

	return nil
}

func getDriver() migrate.Driver {
	switch databaseURL.Scheme { // TODO: do not build all drivers
	case postgres.Dialect.String():
		driver, err := postgres.NewDriver(context.TODO(), sourceDir, databaseURL)
		if err != nil {
			log.Fatal(err)
		}
		return driver
	}
	log.Fatalf("unknown driver \"%s\"", databaseURL.Scheme)
	return nil
}

func main() {
	switch command.(type) {
	case *createCommand:
		e := migrate.Entry{
			Name:      command.(*createCommand).Name,
			Timestamp: now(),
		}
		if err := e.Create(sourceDir); err != nil {
			log.Fatal(err)
		}
		return
	case *migrateCommand:
		driver := getDriver()
		if err := driver.EnsureMigrationsTable(context.Background(), historyTable); err != nil {
			log.Fatal(err)
		}

		count := command.(*migrateCommand).Count
		if count > 0 {
			if err := up(driver, count); err != nil {
				log.Fatal(err)
			}
		} else if count < 0 {
			if err := down(driver, -count); err != nil {
				log.Fatal(err)
			}
		} else {
			panic("count == 0")
		}
	case *testCommand:
		driver := getDriver()
		ok, err := driver.Ping(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		if ok {
			os.Exit(0)
		} else {
			// see /usr/include/sysexits.h EX_UNAVAILABLE
			os.Exit(69)
		}
	default:
		log.Fatal("unknown command")
	}
}
