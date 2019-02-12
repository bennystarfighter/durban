package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/fatih/color"
	"github.com/spf13/viper"
)

var er = color.New(color.BgRed).Add(color.FgBlack)
var suc = color.New(color.BgGreen).Add(color.FgBlack)

func main() {
	fmt.Println("Starting program!")
	sc := make(chan os.Signal, 1)
	go CLIhandler(&sc)
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	if runtime.GOOS == "linux" {
		viper.AddConfigPath("$HOME/.config/durband/")
		viper.AddConfigPath(".")
	}

	err := viper.ReadInConfig()
	if err != nil {
		er.Println("Error reading config file!", err.Error())
		return
	}

	token := viper.GetString("token")

	discord, err := discordgo.New("Bot " + token)
	if err != nil {
		er.Println("Error starting bot!")
	}

	err = discord.Open()
	if err != nil {
		er.Println(("error opening connection to discord,"), err)
		return
	}
	go discord.AddHandler(messageCreate)
	suc.Println("Bot is now running.  Press CTRL-C to exit or write exit and press enter!")
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	discord.Close()
}

func CLIhandler(sc *chan os.Signal) {
	for {
		var input string
		fmt.Scanf("%s", &input)
		if input == "stop" {
			*sc <- os.Kill
		}
	}
}
