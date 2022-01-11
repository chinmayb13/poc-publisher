package config

import (
	"github.com/spf13/viper"
)

//stores all conf values read by viper from config env files

func setupViperConfig() {
	viper.AddConfigPath("../../config") //for debugging
	viper.AddConfigPath("./config")     //for binary
	viper.AddConfigPath("./app/config") //for docker
	viper.SetConfigName("config.local")
}

func init() {
	setupViperConfig()
}

type AppConfig struct {
	DB     DBConfig
	PubSub PubSubCfg
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

	err = viper.Unmarshal(&db)
	if err != nil {
		return
	}

	err = viper.Unmarshal(&pubsub)
	if err != nil {
		return
	}
	config.DB = db
	config.PubSub = pubsub
	return
}
