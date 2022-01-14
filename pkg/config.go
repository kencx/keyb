package pkg

import (
	"fmt"

	"github.com/spf13/viper"
)

type App struct {
	Prefix   string
	KeyBinds []KeyBind `mapstructure:"keybinds"`
}

type KeyBind struct {
	Command string
	Key     string
}

func GetConfig() map[string]App {

	var config map[string]App
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
