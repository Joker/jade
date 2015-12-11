package main

import (
	"fmt"
	"github.com/Joker/jade"
	"html/template"
	"io/ioutil"
	"net/http"
)

type Person struct {
	Name   string
	Age    int
	Emails []string
	Jobs   []*Job
}

type Job struct {
	Employer string
	Role     string
}

func handler(w http.ResponseWriter, r *http.Request) {

	dat, err := ioutil.ReadFile("template.jade")
	if err != nil {
		fmt.Printf("\nReadFile error: %v", err)
		return
	}
	tmpl, err := jade.Parse("jade_tpl", string(dat))
	if err != nil {
		fmt.Printf("\nParse error: %v", err)
		return
	}

	fmt.Printf("%s", tmpl)

	job1 := Job{Employer: "Monash", Role: "Honorary"}
	job2 := Job{Employer: "Box Hill", Role: "Head of HE"}

	person := Person{
		Name:   "jan",
		Age:    50,
		Emails: []string{"jan@newmarch.name", "jan.newmarch@gmail.com"},
		Jobs:   []*Job{&job1, &job2},
	}

	t, err := template.New("html").Parse(tmpl.String())
	if err != nil {
		fmt.Printf("\nTemplate parse error: %v", err)
		return
	}

	err = t.Execute(w, person)
	if err != nil {
		fmt.Printf("\nExecute error: %v", err)
		return
	}
}

func main() {
	fmt.Println("open  http://localhost:8080/")
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
