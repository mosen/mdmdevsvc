package main

import (
	"database/sql"
	"fmt"
	"github.com/DavidHuie/gomigrate"
	"github.com/containous/flaeg"
	"github.com/containous/staert"
	"github.com/go-kit/kit/log"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/micromdm/dep"
	"github.com/mosen/devicestore/depsync"
	"github.com/mosen/devicestore/device"
	"golang.org/x/net/context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type DEPInfo struct {
	URL            string `description:"DEP url"`
	ConsumerKey    string `description:"Consumer key" short:"k"`
	ConsumerSecret string `description:"Consumer secret" short:"s"`
	AccessToken    string `description:"Access token" short:"t"`
	AccessSecret   string `description:"Access secret" short:"a"`
	SyncInterval   int    `description:"Sync interval (in seconds)" short:"i"`
}

type DatabaseInfo struct {
	Host     string `description:"Hostname or IP address of postgresql server"`
	Port     string `description:"database port number"`
	Name     string `description:"database name"`
	Username string `description:"database username"`
	Password string `description:"database password"`
	SSLMode  string `description:"postgres SSL mode"`
}

type ListenInfo struct {
	IP   string `description:"IP Address to listen on"`
	Port string `description:"listen on port number"`
}

type Configuration struct {
	Db     *DatabaseInfo `description:"Database connection options"`
	Listen *ListenInfo   `description:"Listen"`
	Dep    *DEPInfo      `description:"Device enrollment program options"`
}

func main() {
	var config *Configuration = &Configuration{
		&DatabaseInfo{
			Host:     "localhost",
			Port:     "5432",
			Name:     "mdmdevsvc",
			Username: "mdmdevsvc",
			Password: "mdmdevsvc",
			SSLMode:  "disable",
		},
		&ListenInfo{
			IP:   "0.0.0.0",
			Port: "8080",
		},
		&DEPInfo{
			AccessSecret:   "AS_c31afd7a09691d83548489336e8ff1cb11b82b6bca13f793344496a556b1f4972eaff4dde6deb5ac9cf076fdfa97ec97699c34d515947b9cf9ed31c99dded6ba",
			AccessToken:    "AT_927696831c59ba510cfe4ec1a69e5267c19881257d4bca2906a99d0785b785a6f6fdeb09774954fdd5e2d0ad952e3af52c6d8d2f21c924ba0caf4a031c158b89",
			ConsumerKey:    "CK_48dd68d198350f51258e885ce9a5c37ab7f98543c4a697323d75682a6c10a32501cb247e3db08105db868f73f2c972bdb6ae77112aea803b9219eb52689d42e6",
			ConsumerSecret: "CS_34c7b2b531a600d99a0e4edcf4a78ded79b86ef318118c2f5bcfee1b011108c32d5302df801adbe29d446eb78f02b13144e323eb9aad51c79f01e50cb45c3a68",
			SyncInterval:   60,
		},
	}

	var pointersConfig *Configuration = &Configuration{}

	rootCmd := &flaeg.Command{
		Name:                  "mdmdevsvc",
		Description:           "MDM device service stores information about devices under management.",
		Config:                config,
		DefaultPointersConfig: pointersConfig,
		Run: func() error {
			run(config)
			return nil
		},
	}

	st := staert.NewStaert(rootCmd)
	toml := staert.NewTomlSource("mdmdevsvc", []string{"./"})
	fl := flaeg.New(rootCmd, os.Args[1:])

	st.AddSource(toml)
	st.AddSource(fl)
	if _, err := st.LoadConfig(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	if err := st.Run(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	os.Exit(0)
}

func run(config *Configuration) {
	var err error
	var db *sql.DB
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.NewContext(logger).With("ts", log.DefaultTimestampUTC)
	}

	db, err = sql.Open("postgres", fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Db.Host,
		config.Db.Port,
		config.Db.Username,
		config.Db.Password,
		config.Db.Name,
	))
	if err != nil {
		logger.Log("level", "error", "msg", err)
		os.Exit(-1)
	}

	err = db.Ping()
	if err != nil {
		logger.Log("level", "error", "msg", err)
		os.Exit(-1)
	}

	var dbx *sqlx.DB = sqlx.NewDb(db, "postgres")

	migrator, _ := gomigrate.NewMigrator(db, gomigrate.Postgres{}, "./migrations")
	migrationErr := migrator.Migrate()

	if migrationErr != nil {
		logger.Log("level", "error", "msg", err)
		os.Exit(-1)
	}

	ctx := context.Background()

	deviceDb := device.NewRepository(dbx)
	deviceSvc := device.NewService(deviceDb)
	deviceSvc = device.LoggingMiddleware(log.NewContext(logger).With("component", "device.Service"))(deviceSvc)
	deviceHandler := device.MakeHTTPHandler(ctx, deviceSvc, logger)

	mux := http.NewServeMux()

	mux.Handle("/v1/devices/", deviceHandler)

	portStr := fmt.Sprintf("%v:%v", config.Listen.IP, config.Listen.Port)

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		logger.Log("level", "info", "msg", "Listening on "+portStr)
		errs <- http.ListenAndServe(portStr, mux)
	}()

	depTicker := time.NewTicker(time.Duration(config.Dep.SyncInterval) * time.Second)

	depConfig := &dep.Config{
		ConsumerKey:    config.Dep.ConsumerKey,
		ConsumerSecret: config.Dep.ConsumerSecret,
		AccessSecret:   config.Dep.AccessSecret,
		AccessToken:    config.Dep.AccessToken,
	}

	depClient, err := dep.NewClient(depConfig, dep.ServerURL(config.Dep.URL))
	syncer := depsync.NewSyncer(depClient, log.NewContext(logger).With("component", "depsync.Syncer"), depTicker.C)
	depDeviceChan := make(chan dep.Device)

	depWriter := depsync.NewWriter(deviceDb, log.NewContext(logger).With("component", "depsync.Writer"))

	go syncer.Start(depDeviceChan)
	go depWriter.Start(depDeviceChan)

	logger.Log("exit", <-errs)
}
