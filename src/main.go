package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	db "Athena-Backend/database"
)

func colorize(text string, colorCode string) string {
	reset := "\033[0m"
	return colorCode + text + reset
}

func main() {
	port := 3551
	addr := "127.0.0.1:" + strconv.Itoa(port)
	mongoURI := "mongodb://127.0.0.1/AthenaBackend"

	brightCyan := "\033[96m"
	tag := colorize("[Backend]", brightCyan)

	client, err := db.ConnectMongo(mongoURI, tag)
	if err != nil {
		panic(fmt.Sprintf("%s Failed to connect to MongoDB: %v", tag, err))
	}
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			fmt.Printf("%s Error disconnecting MongoDB: %v\n", tag, err)
		}
	}()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Athena Backend Running On Port :3 %d", port)
	})

	fmt.Printf("%s Starting Athena Backend on %s\n", tag, addr)

	if err := http.ListenAndServe(addr, nil); err != nil {
		panic(err)
	}
}
