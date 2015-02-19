package main

import (
	"fmt"
	"io/ioutil"
	"github.com/Joker/jade"
)


func main() {
	dat, err := ioutil.ReadFile("template.jade")
	if err != nil {
		fmt.Printf("ReadFile error: %v", err)
		return
	}

	tmpl, err := jade.New("jade_tpl").Parse(string(dat), "", "", make( map[string]*jade.Tree ), nil)
	if err != nil {
		fmt.Printf("Parse error: %v", err)
		return
	}

	fmt.Printf( "\nOutput:\n\n"  )
	fmt.Printf( tmpl.Root.String() )
}
