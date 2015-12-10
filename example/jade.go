package main

import (
	"fmt"
	"io/ioutil"
	"github.com/joker/jade"
)


func main() {
	dat, err := ioutil.ReadFile("template.jade")
	if err != nil {
		fmt.Printf("ReadFile error: %v", err)
		return
	}

	tmpl, err := jade.Parse("name_of_tpl", string(dat))
	if err != nil {
		fmt.Printf("Parse error: %v", err)
		return
	}

	fmt.Printf( "\nOutput:\n\n%s", tmpl )
}
