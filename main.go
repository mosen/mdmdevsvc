package main

import (

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
	"github.com/containous/flaeg"
	"github.com/containous/staert"
)


type HostInfo struct {
	Host string `description:"Hostname or IP address"`
	Port string `description:"Port number"`
}

type Configuration struct {
	Db HostInfo `description:"Database"`
	Listen HostInfo `description:"Listen"`
}

func main() {
	var config *Configuration = &Configuration{
		Db: HostInfo{
			Host: "localhost",
			Port: "5432",
		},
		Listen: HostInfo{
			Host: "127.0.0.1",
			Port: "8080",
		},
	}

	var pointersConfig *Configuration = &Configuration{}

	rootCmd := &flaeg.Command{
		Name: "mdmdevsvc",
		Description: "MDM device service stores information about devices under management.",
		Config: config,
		DefaultPointersConfig: pointersConfig,
		Run: func() error {
			fmt.Printf("Run flaegtest command with config : %+v\n", config)
			return nil
		},
	}

	st := staert.NewStaert(rootCmd)
	toml := staert.NewTomlSource("mdmappsvc", []string{"./"})
	fl := flaeg.New(rootCmd, os.Args[1:])

	st.AddSource(toml)
	st.AddSource(fl)
	loadedConfig, err := st.LoadConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	ctx := context.Background()
	logger := getLogger()

	db, err := sqlx.Open("postgres", loadedConfig.(Configuration).Db.Host)
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

	portStr := fmt.Sprintf(":%v", loadedConfig.(Configuration).Listen.Port)
	http.ListenAndServe(portStr, nil)
}

func getEnvDefault(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		value = defaultValue
	}
	return value
}

func getLogger() log.Logger {
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

	logger := term.NewLogger(os.Stdout, log.NewJSONLogger, colorFn)
	return logger
}
