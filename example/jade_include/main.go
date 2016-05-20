package main

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/Joker/jade"
)

func ReadAndParse(path string) string {
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("\nReadFile error: %v", err)
	}
	tpl, err := jade.Parse("jade_tp", string(dat))
	if err != nil {
		log.Printf("\nParse error: %v", err)
	}
	return tpl
}

func handler(w http.ResponseWriter, r *http.Request) {
	jade_tpl := ReadAndParse("template.jade")
	log.Printf("%s\n\n", jade_tpl)

	//

	funcMap := template.FuncMap{
		"include": func(includePath string) (template.HTML, error) {
			include_tpl := ReadAndParse(includePath)
			log.Printf("%s\n\n", include_tpl)

			go_partial_tpl, _ := template.New("partial").Parse(include_tpl)

			buf := new(bytes.Buffer)
			go_partial_tpl.Execute(buf, "")
			return template.HTML(buf.String()), nil

		},
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

func main() {
	log.Println("open  http://localhost:8080/")
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
