package config

import (
	"github.com/micro/go-config"
	"github.com/micro/go-config/source/env"
	"github.com/micro/go-config/source/file"
	"log"
	"os"
	"path"
)

type DatabaseConfig struct {
	Driver   string `json:"driver"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
}

type RedisConfig struct {
	Address string `json:"address"`
	Port    string `json:"port"`
}

type TicketConfig struct {
	Driver string `json:"driver"` // ticket 使用的驱动 只支持 redis 和 database
	TTL    int64  `json:"ttl"`    // ticket 的过期时间 （毫秒）
}

type FilesystemConfig struct {
	Driver string `json:"driver"`
	Root   string `json:"root"`
}

type MinioConfig struct {
	Host       string `json:"host"`
	AccessKey  string `json:"accesskey"`
	SecretKey  string `json:"secretkey"`
	SSL        string `json:"ssl"`
	BucketName string `json:"bucketname"`
}

type ImageproxyConfig struct {
	Host        string
	OmitBaseUrl bool `json:"omitbaseurl"`
}

type Config struct {
	EnvVarPrefix string           `json:"env-var-prefix"`
	ServiceName  string           `json:"service-name"`
	ServerAddr   string           `json:"server-addr"` // addr:port
	AppSalt      string           `json:"appsalt"`
	QueueNum     int              `json:"queue-num"`
	Fs           FilesystemConfig `json:"filesystem"`
	DB           DatabaseConfig   `json:"database"`
	Redis        RedisConfig      `json:"redis"`
	Ticket       TicketConfig     `json:"ticket"`
	Minio        MinioConfig      `json:"minio"`
	Imageproxy   ImageproxyConfig
}

func LoadConfig(filepath string) *Config {
	c := &Config{}
	pwd, _ := os.Getwd()
	fileSource := file.NewSource(file.WithPath(path.Join(pwd, filepath)))
	checkErr(config.Load(fileSource))
	// env 的配置会覆盖文件中的配置
	envSource := env.NewSource(env.WithStrippedPrefix(config.Get("env-var-prefix").String("CLOUD")))
	checkErr(config.Load(envSource))
	checkErr(config.Scan(c))
	return c
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
