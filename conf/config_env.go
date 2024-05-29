package conf

import (
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

// AppConfig presents app conf
type AppConfig struct {
	ServerEnv string `env:"SERVER_ENV" envDefault:"stg"`
	Port      string `env:"PORT"`
	// Database
	DBHost        string `env:"DB_HOST"`
	DBPort        string `env:"DB_PORT"`
	DBUser        string `env:"DB_USER"`
	DBPass        string `env:"DB_PASS"`
	DBName        string `env:"DB_NAME"`
	EnableDB      string `env:"ENABLE_DB" envDefault:"true"`
	DbDebugEnable bool   `env:"DB_DEBUG_ENABLE" envDefault:"true"`

	// Media
	AWSUpCloudRegion      string `env:"UC_AWS_REGION"`
	AWSMediaDomain        string `env:"UC_AWS_MEDIA_DOMAIN"`
	AWSMediaFullDomain    string `env:"UC_AWS_MEDIA_FULL_DOMAIN" envDefault:"https://noormatch.ap-south-1.linodeobjects.com/"`
	AWSUpCloudAccessKeyID string `env:"UC_AWS_ACCESS_KEY_ID"`
	AWSUpCloudSecretKey   string `env:"UC_AWS_SECRET_KEY"`
	AWSBucket             string `env:"UC_AWS_BUCKET"`
	AWSS3ACL              string `env:"UC_AWS_S3_ACL"`
	MediaBaseProxyURL     string `env:"MEDIA_BASE_PROXY_URL"`
}

var config AppConfig

func SetEnv() {
	err := godotenv.Load("./.env")
	if err != nil {
		fmt.Println(err)
	}
	_ = env.Parse(&config)
}

func LoadEnv() AppConfig {
	return config
}
