package config

import (
	"io/ioutil"
	"log"
	"strings"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Database struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Port     string `yaml:"port"`
		IP       string `yaml:"ip"`
		DBName   string `yaml:"dbname"`
		Charset  string `yaml:"charset"`
	}
	SSH struct {
		Addr string `yaml:"addr"`
		Pswd string `yaml:"pswd"`
	}
}

var Conf = &Config{}

func MustLoadConfig(path string) {
	File, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("read config file error: #%v", err)
	}
	err = yaml.Unmarshal(File, Conf)
	if err != nil {
		log.Fatalf("config decodes error: %v", err)
	}
}

func GetDsn() string {
	username := Conf.Database.Username
	password := Conf.Database.Password
	ip := Conf.Database.IP
	port := Conf.Database.Port
	dbname := Conf.Database.DBName
	charset := Conf.Database.Charset
	return strings.Join([]string{username, ":", password, "@mysql+tcp(", ip, ":", port, ")/", dbname, "?charset=", charset}, "")
}
