package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func colorize(text string, colorCode string) string {
	reset := "\033[0m"
	return colorCode + text + reset
}

func main() {
	port := 3551
	addr := "127.0.0.1:" + strconv.Itoa(port)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Athena Backend Running On Port :3 %d", port)
	})
	brightCyan := "\033[96m"
	tag := colorize("[Backend]", brightCyan)

	fmt.Printf("%s Starting Athena Backend on %s\n", tag, addr)

	if err := http.ListenAndServe(addr, nil); err != nil {
		panic(err)
	}
}
