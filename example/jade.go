package main

import (
	"fmt"
	"io/ioutil"

	"github.com/Joker/hpp"
	"github.com/Joker/jade"
)

func main() {
	dat, err := ioutil.ReadFile("template.jade")
	if err != nil {
		fmt.Printf("ReadFile error: %v", err)
		return
	}

	tmpl, err := jade.Parse("name_of_tpl", dat)
	if err != nil {
		fmt.Printf("Parse error: %v", err)
		return
	}

	fmt.Printf("\nOutput:\n\n%s", hpp.PrPrint(tmpl))
}
