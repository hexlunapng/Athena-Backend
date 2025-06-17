package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func main() {
	port := 3551
	addr := ":" + strconv.Itoa(port)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Athena Backend Running On Port :3 %d", port)
	})

	fmt.Printf("Starting Athena Backend on %s\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		panic(err)
	}
}
