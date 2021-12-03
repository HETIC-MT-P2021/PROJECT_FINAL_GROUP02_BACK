package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/SteakBarbare/RPGBot/database"
	"github.com/SteakBarbare/RPGBot/handlers"
	"github.com/bwmarrin/discordgo"
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"github.com/SteakBarbare/RPGBot/socketio"
)

func main() {

	var token string

	// Loading env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	token = os.Getenv("TOKEN")

	cfg := database.Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err)
	}

	database.Connect(cfg)

	// Start socketio connection
	go func() { socketio.Connect() }()

	// Create a new Discord session using the provided bot token.

	dg, err := discordgo.New(`Bot ` + token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(handlers.MessageCreate)
	dg.AddHandler(handlers.ReactionsHandler)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}
