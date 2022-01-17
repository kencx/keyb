package main

import (
	"fmt"

	"github.com/spf13/viper"
)

type Program struct {
	Prefix   string    `mapstructure:"prefix,omitempty"`
	KeyBinds []KeyBind `mapstructure:"keybinds"`
}

type KeyBind struct {
	Command string
	Key     string
	Comment string `mapstructure:"comment,omitempty"`
}

func GetConfig() map[string]Program {

	var config map[string]Program
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("config file not found")
		}
		fmt.Printf("error in config: %v", err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		fmt.Printf("failed to decode: %v", err)
	}
	return config
}
