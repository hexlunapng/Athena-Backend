package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/joho/godotenv"

	db "Athena-Backend/database"
	bot "Athena-Backend/discord"
)

const (
	brightCyan = "\033[96m"
	blurple    = "\033[38;5;63m"
	resetColor = "\033[0m"
)

func colorize(text string, colorCode string) string {
	return colorCode + text + resetColor
}

func colorizeBackend(text string) string {
	return colorize(text, brightCyan)
}

func colorizeDiscord(text string) string {
	return colorize(text, blurple)
}

func main() {
	_ = godotenv.Load()

	port := 3551
	addr := "127.0.0.1:" + strconv.Itoa(port)
	mongoURI := "mongodb://127.0.0.1/AthenaBackend"
	discordToken := os.Getenv("DISCORD_BOT_TOKEN")

	tagBackend := colorizeBackend("[BACKEND]")
	tagDiscord := colorizeDiscord("[DISCORD]")

	client, err := db.ConnectMongo(mongoURI, tagBackend)
	if err != nil {
		panic(fmt.Sprintf("%s Failed to connect to MongoDB: %v", tagBackend, err))
	}
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			fmt.Printf("%s Error disconnecting MongoDB: %v\n", tagBackend, err)
		}
	}()

	dg, err := bot.StartAthenaBackendDiscordBot(discordToken)
	if err != nil {
		panic(fmt.Sprintf("%s Discord bot error: %v", tagDiscord, err))
	}
	defer func() {
		dg.Close()
		fmt.Println(colorizeDiscord("[DISCORD] Bot shut down."))
	}()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Athena Backend Running On Port :3 %d", port)
	})
	fmt.Printf("%s Starting Athena Backend on %s\n", tagBackend, addr)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		if err := http.ListenAndServe(addr, nil); err != nil {
			panic(err)
		}
	}()

	<-stop
	fmt.Println(colorizeBackend("[BACKEND] Shutting down everything."))
}
