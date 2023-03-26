package main

import (
	"embed"
	"fmt"
	"github.com/JonasGruenwald/mp3/cmd"
	"github.com/spf13/viper"
	"os"
)

//go:embed templates/*
var templateFs embed.FS

func main() {
	// Set up config
	viper.SetConfigName("mp3-config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME")
	viper.SetDefault("AdoptedServices", []string{})
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error
		} else {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	// Run command handler
	cmd.Execute(templateFs)
}
