package config

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/spf13/viper"
)

//stores all conf values read by viper from config env files

func setupViperConfig() {
	configPath := os.Getenv("APP_DIR") + "/config"
	viper.AddConfigPath(configPath) //for debugging
	//viper.AddConfigPath("./config")     //for binary
	//viper.AddConfigPath("./app/config") //for docker
	viper.SetConfigName("config.local")
}

func setRootIfNotExist() {
	_, present := os.LookupEnv("APP_DIR")
	if !present {
		_, b, _, _ := runtime.Caller(0)
		os.Setenv("APP_DIR", filepath.Dir(filepath.Dir(b)))
	}
}

func init() {
	setRootIfNotExist()
	setupViperConfig()
}

type AppConfig struct {
	DB           DBConfig
	PubSub       PubSubCfg
	DeployConfig DeployConfig
}

type DeployConfig struct {
	Port   string `mapstructure:"PORT"`
	AppDir string `mapstructure:"APP_DIR"`
}

type DBConfig struct {
	Host      string `mapstructure:"DB_HOST"`
	Port      int    `mapstructure:"DB_PORT"`
	QueueSize int    `mapstructure:"DB_QUEUE_SIZE"`
	LimitConn bool   `mapstructure:"DB_LIMIT_CONN"`
	Timeout   int    `mapstructure:"DB_TIME_OUT"`
	Set       string `mapstructure:"DB_SET"`
	NameSpace string `mapstructure:"DB_NAME_SPACE"`
}

type PubSubCfg struct {
	ProjectID string `mapstructure:"PUBSUB_PROJECT_ID"`
	TopicID   string `mapstructure:"PUBSUB_TOPIC_ID"`
}

func LoadConfig(path string) (config AppConfig, err error) {
	//viper.AddConfigPath(path)
	//viper.SetConfigName("config.local")
	//viper.SetConfigType("env")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	var db DBConfig
	var pubsub PubSubCfg
	var deploy DeployConfig

	err = viper.Unmarshal(&db)
	if err != nil {
		return
	}

	err = viper.Unmarshal(&pubsub)
	if err != nil {
		return
	}

	err = viper.Unmarshal(&deploy)
	if err != nil {
		return
	}

	config.DB = db
	config.PubSub = pubsub
	config.DeployConfig = deploy
	return
}
