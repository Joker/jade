package jade

/*
Parse parses the template definition string to construct a representation of the template for execution.

Trivial usage:

	package main

	import (
		"fmt"
		"github.com/joker/jade"
	)

	func main() {
		tpl, err := jade.Parse("name_of_tpl", "doctype 5: html: body: p Hello world!")
		if err != nil {
			fmt.Printf("Parse error: %v", err)
			return
		}

		fmt.Printf( "\nOutput:\n\n%s", tpl  )
	}

Output:

	<!DOCTYPE html>
	<html>
	    <body>
	        <p>Hello world!</p>
	    </body>
	</html>
*/
func Parse(name, text string) (string, error) {
	outTpl, err := newTree(name).Parse(text, leftDelim, rightDelim, make(map[string]*tree))
	return outTpl.String(), err
}

func (t *tree) String() string {
	return t.Root.String()
}
