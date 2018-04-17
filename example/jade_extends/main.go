package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/Joker/hpp"
	"github.com/Joker/jade"
)

func handler(w http.ResponseWriter, r *http.Request) {
	index, err := jade.ParseFile("index.jade")
	if err != nil {
		log.Printf("\nParseFile error: %v", err)
	}
	log.Printf("%s\n\n", hpp.PrPrint(index))

	//

	go_tpl, err := template.New("layout").Parse(index)
	if err != nil {
		log.Printf("\nTemplate parse error: %v", err)
	}

	err = go_tpl.Execute(w, "")
	if err != nil {
		log.Printf("\nExecute error: %v", err)
	}
}

func main() {
	log.Println("open  http://localhost:8080/")
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
