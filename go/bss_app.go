package main

import (
	"bss/go/data"
	"bss/go/handlers"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	log.Println("Initializing HTTP Server")

	index := handlers.BuildPageHandler(
		"index",
		"./react/build",
		"",
		"static/css",
		"static/js",
	)

	data, err := data.BuildData(
		"sqlite3",
		"./db.db",
	)
	if err != nil {
		log.Printf(
			"Error Building Data Object\n\t:%s",
			err,
		)
		return
	}

	itemHandler := handlers.ItemHandler{
		Data: data,
	}

	http.HandleFunc("/", index.Handle)
	http.HandleFunc("/items/update", itemHandler.HandleUpdate)
	http.HandleFunc("/items/query", itemHandler.HandleQuery)

	http.ListenAndServe(":80", nil)
}
