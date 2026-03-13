package config

import "github.com/spf13/viper"

type Config struct {
	CommonEtcd struct {
		Host string `mapstructure:"host"`
	} `mapstructure:"common-etcd"`
	CommonRedis struct {
		RedisAddress  string `mapstructure:"redisAddress"`
		RedisPassword string `mapstructure:"redisPassword"`
		DB            int    `mapstructure:"db"`
	} `mapstructure:"common-redis"`
	CommonRabbitMQ struct {
		URL       string `mapstructure:"url"`
		QueueName string `mapstructure:"queueName"`
	} `mapstructure:"common-rabbitmq"`
	CommonMysql struct {
		Dsn string `mapstructure:"dsn"`
	} `mapstructure:"common-mysql"`
}

var AppConfig Config

// Init 初始化配置
func Init() error {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath("./config/")

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	return viper.Unmarshal(&AppConfig)
}
