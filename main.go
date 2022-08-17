package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/OleksiiKhanin/companysvc/api"
	"github.com/OleksiiKhanin/companysvc/config"
	"github.com/OleksiiKhanin/companysvc/db"
	"github.com/OleksiiKhanin/companysvc/domain"
	"github.com/OleksiiKhanin/companysvc/service"

	"github.com/gorilla/mux"
	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

func initConfig() (*config.Config, error) {
	var conf config.Config
	confPath := os.Getenv("CONFIG")
	if confPath == "" {
		conf, _ := yaml.Marshal(conf)
		return nil, fmt.Errorf(
			"Please specified a configuration file path in the CONFIG environment variable\n The config.yaml example:\n%s",
			conf,
		)
	}
	viper.AddConfigPath(".")
	viper.AddConfigPath("/")
	viper.SetConfigName(confPath)
	viper.SetEnvPrefix("app")                                        // You can use environment variable with the same name as config file and prefix APP_
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_")) // '.' -> '_' and '-' -> '_' in env variable
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config file: %w", err)
	}
	if err := viper.Unmarshal(&conf); err != nil {
		return nil, fmt.Errorf("parse config file %w", err)
	}
	return &conf, nil
}

func initDB(c *config.DatabaseConfig) (*sql.DB, error) {
	connectionString := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.URL, c.Port, c.Login, c.Password, c.NameDB,
	)
	storage, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("try open db connection: %w", err)
	}
	if c.MaxConns > 0 {
		storage.SetMaxOpenConns(c.MaxConns)
	}
	if err := storage.Ping(); err != nil {
		return nil, fmt.Errorf("ping connection to db: %w", err)
	}
	return storage, err
}

func initQueue(c *config.QueueConfig) (*nats.Conn, error) {
	conn, err := nats.Connect(
		fmt.Sprintf("%s:%d", c.URL, c.Port),
		nats.RetryOnFailedConnect(true),
		nats.ReconnectWait(c.ReconnectWait),
		nats.PingInterval(c.PingInterval),
	)
	if err != nil {
		return nil, fmt.Errorf("create connection to queue: %w", err)
	}
	return conn, nil
}

func startServer(c *config.ServerConfig, iCompany domain.ICompany) {
	r := mux.NewRouter().UseEncodedPath()
	api.InitAPI(
		r.PathPrefix(c.PrefixAPI).Subrouter(),
		iCompany,
		log.StandardLogger(),
	)
	log.StandardLogger().Infof("Start listening at %s", c.URL)
	log.Fatal(http.ListenAndServe(c.URL, r))
}

func main() {
	c, err := initConfig()
	if err != nil {
		log.Fatalln(err.Error())
		os.Exit(1)
	}

	// Set log level if config.logLevel correct, otherwise use default logLevel (INFO)
	if logLevel, err := log.ParseLevel(c.LogLevel); err == nil {
		log.SetLevel(logLevel)
	}

	storage, err := initDB(&c.Db)
	if err != nil {
		log.Fatalln(err.Error())
		os.Exit(1)
	}
	defer storage.Close()

	if c.Db.Migrations != "" {
		if err := db.MigrateSchema(storage, c.Db.Migrations); err != nil {
			log.Fatalln(err.Error())
			os.Exit(1)
		}
	}

	queue, err := initQueue(&c.Event)
	if err != nil {
		log.Error(err.Error()) // log error but continue
	} else {
		defer queue.Close()
	}

	iCompany := service.NewCompanyService(
		db.NewCompanyPostgresRepo(storage, log.StandardLogger()),
		queue,
		service.GetResolverIPAPI(c.Loc.URL, c.Loc.RetryAttempt, log.StandardLogger()),
		c.Event.EventChannel,
		log.StandardLogger(),
		c.Loc.AllowedCountiesCodes...,
	)

	startServer(&c.Server, iCompany)
}
