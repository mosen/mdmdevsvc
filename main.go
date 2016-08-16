package main

import (
	"flag"
	"fmt"
	"github.com/DavidHuie/gomigrate"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/term"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/mosen/devicestore/device"
	"golang.org/x/net/context"
	"net/http"
	"os"
)

func main() {
	colorFn := func(keyvals ...interface{}) term.FgBgColor {
		for i := 0; i < len(keyvals)-1; i += 2 {
			if keyvals[i] != "level" {
				continue
			}
			switch keyvals[i+1] {
			case "debug":
				return term.FgBgColor{Fg: term.DarkGray}
			case "info":
				return term.FgBgColor{Fg: term.Gray}
			case "warn":
				return term.FgBgColor{Fg: term.Yellow}
			case "error":
				return term.FgBgColor{Fg: term.Red}
			case "crit":
				return term.FgBgColor{Fg: term.Gray, Bg: term.DarkRed}
			default:
				return term.FgBgColor{}
			}
		}
		return term.FgBgColor{}
	}

	ctx := context.Background()
	logger := term.NewLogger(os.Stdout, log.NewJSONLogger, colorFn)

	// Flags
	var (
		flPort      = flag.String("port", getEnvDefault("DEVICESVC_PORT", "3000"), "port to listen on")
		flDbConnUrl = flag.String("db", getEnvDefault("DEVICESVC_DB_CONN", "user=devicestore password=devicestore dbname=devicestore sslmode=disable"), "database connection url (postgres)")
	)
	flag.Parse()

	var db *sqlx.DB
	var err error
	db, err = sqlx.Open("postgres", *flDbConnUrl)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = db.Ping()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	migrator, _ := gomigrate.NewMigrator(db.DB, gomigrate.Postgres{}, "./migrations")
	migrationErr := migrator.Migrate()

	if migrationErr != nil {
		logger.Log("err", err)
		os.Exit(1)
	}

	deviceDb := device.NewRepository(db)
	deviceSvc := device.NewService(deviceDb)
	deviceSvc = device.LoggingMiddleware(logger)(deviceSvc)
	deviceHandler := device.MakeHTTPHandler(ctx, deviceSvc, logger)

	mux := http.NewServeMux()

	mux.Handle("/v1/", deviceHandler)

	portStr := fmt.Sprintf(":%v", flPort)
	http.ListenAndServe(portStr, nil)
}

func getEnvDefault(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		value = defaultValue
	}
	return value
}
