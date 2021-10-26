package main

import (
	"log"
	"net/http"
)

func main() {
	foo := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	}

	http.HandleFunc("/foo", foo)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
