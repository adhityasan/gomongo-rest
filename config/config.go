package config

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

type tomlConfig struct {
	Owner     ownerInfo
	Mongo     mongoConnection `toml:"mongoconnection"`
	DBModules map[string]dbmodule
	Azure     azure
	Aws       aws
	App       app
	Dukcapil  dukcapil
}

type app struct {
	Name string
	Host string `toml:"hostname"`
	Port string
	Src  string
}

type ownerInfo struct {
	Name string
	Org  string `toml:"organization"`
}

type mongoConnection struct {
	Host     string `toml:"hostname"`
	Username string
	Port     string
	Password string
	URL      string
}

type dbmodule struct {
	Db   string `toml:"database"`
	Coll string `toml:"collection"`
}

type dukcapil struct {
	Endpoint string
	key      string
}

type aws struct {
	Region    string // "us-east-1" / "ap-northeast-2"
	KeyID     string
	SecretKey string
}

type azure struct {
	Endpoint string
	APIKey   string `toml:"key"`
}

// Of save all the configuration in Config.toml
var Of tomlConfig

// SetConfiguration asdasdasdasd
func init() {
	fmt.Println("Init config")

	if _, err := toml.DecodeFile("Config.toml", &Of); err != nil {
		fmt.Printf("Load Config.toml Failure : %v\n", err)
		return
	}

	for dbmodname, dbmodule := range Of.DBModules {
		fmt.Printf("DBModules: %s (%s, %s)\n", dbmodname, dbmodule.Db, dbmodule.Coll)
	}
}
