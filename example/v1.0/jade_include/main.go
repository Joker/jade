package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/Joker/hpp"
	"github.com/Joker/jade"
)

func handler(w http.ResponseWriter, r *http.Request) {
	jade_tpl, err := jade.ParseFile("template.jade")
	if err != nil {
		log.Printf("\nParseFile error: %v", err)
	}
	log.Printf("%s\n\n", hpp.PrPrint(jade_tpl))

	//

	funcMap := template.FuncMap{
		"bold": func(content string) (template.HTML, error) {
			return template.HTML("<b>" + content + "</b>"), nil
		},
	}

	//

	go_tpl, err := template.New("html").Funcs(funcMap).Parse(jade_tpl)
	if err != nil {
		log.Printf("\nTemplate parse error: %v", err)
	}

	err = go_tpl.Execute(w, "")
	if err != nil {
		log.Printf("\nExecute error: %v", err)
	}
}

func js(w http.ResponseWriter, r *http.Request) {}

func main() {
	log.Println("open  http://localhost:8080/")
	http.HandleFunc("/javascripts/", js)
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
