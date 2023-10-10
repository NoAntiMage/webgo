package config

import (
	"encoding/json"
	"fmt"
	"goweb/common/jsonx"
	"goweb/common/validx"
	"os"
)

var globalConf GlobalConfig

type GlobalConfig struct {
	LogConf LogConfig `json:"log"`
	DbConf  DbConfig  `json:"db"`
}

func GetGlobalConfig() GlobalConfig {
	return globalConf
}

func ConfigInit() {
	fmt.Println("config init...")
	//CUSTOM: define your configuration here...
	fileBytes, err := os.ReadFile("./configuration-example.json")
	// fileBytes, err := os.ReadFile("./configuration.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(fileBytes, &globalConf)
	if err != nil {
		panic(err)
	}
	if err = jsonx.Struct2StructWithRule(&globalConf, jsonx.RuleDefault); err != nil {
		panic(err)
	}

	if err = validx.GetValidator().Struct(globalConf); err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", globalConf)
}
