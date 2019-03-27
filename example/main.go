package main

//go:generate jade -pkg=main -writer -fmt index.jade

import (
	"net/http"
)

func main() {

	http.HandleFunc("/", func(wr http.ResponseWriter, req *http.Request) {
		Index("Jade.go", true, wr)
	})

	http.ListenAndServe(":8080", nil)
}
