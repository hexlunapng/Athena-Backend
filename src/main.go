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

func colorize(text string, colorCode string) string {
	reset := "\033[0m"
	return colorCode + text + reset
}

func main() {
	_ = godotenv.Load()

	port := 3551
	addr := "127.0.0.1:" + strconv.Itoa(port)
	mongoURI := "mongodb://127.0.0.1/AthenaBackend"
	discordToken := os.Getenv("DISCORD_BOT_TOKEN")

	brightCyan := "\033[96m"
	tag := colorize("[BACKEND]", brightCyan)

	client, err := db.ConnectMongo(mongoURI, tag)
	if err != nil {
		panic(fmt.Sprintf("%s Failed to connect to MongoDB: %v", tag, err))
	}
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			fmt.Printf("%s Error disconnecting MongoDB: %v\n", tag, err)
		}
	}()

	dg, err := bot.StartAthenaBackendDiscordBot(discordToken)
	if err != nil {
		panic(fmt.Sprintf("%s Discord bot error: %v", tag, err))
	}
	defer func() {
		dg.Close()
		fmt.Println("[DISCORD] Bot shut down.")
	}()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Athena Backend Running On Port :3 %d", port)
	})
	fmt.Printf("%s Starting Athena Backend on %s\n", tag, addr)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		if err := http.ListenAndServe(addr, nil); err != nil {
			panic(err)
		}
	}()

	<-stop
	fmt.Println("[BACKEND] Shutting down everything.")
}
