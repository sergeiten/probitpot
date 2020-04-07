package main

import (
	"log"
	"probitpot"

	"github.com/jessevdk/go-flags"
)

func main() {
	var opts probitpot.Opts

	_, err := flags.Parse(&opts)

	if err != nil {
		log.Fatalf("failed to parse opts: %v", err)
	}

	bot, err := probitpot.NewBot(opts)
	if err != nil {
		log.Fatalf("failed to create bot: %v", err)
	}

	err = bot.Run()
	if err != nil {
		log.Fatalf("failed to run: %v", err)
	}
}
