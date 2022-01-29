// Jade.go - template engine. Package implements Jade-lang templates for generating Go html/template output.
package jade

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"path/filepath"
)

/*
Parse parses the template definition string to construct a representation of the template for execution.

Trivial usage:

	package main

	import (
		"fmt"
		"html/template"
		"net/http"

		"github.com/Joker/jade"
	)

	func handler(w http.ResponseWriter, r *http.Request) {
		jadeTpl, _ := jade.Parse("jade", []byte("doctype 5\n html: body: p Hello #{.Word}!"))
		goTpl, _ := template.New("html").Parse(jadeTpl)

		goTpl.Execute(w, struct{ Word string }{"jade"})
	}

	func main() {
		http.HandleFunc("/", handler)
		http.ListenAndServe(":8080", nil)
	}

Output:

	<!DOCTYPE html><html><body><p>Hello jade!</p></body></html>
*/
func Parse(name string, text []byte) (string, error) {
	outTpl, err := New(name).Parse(text)
	if err != nil {
		return "", err
	}
	b := new(bytes.Buffer)
	outTpl.WriteIn(b)
	return b.String(), nil
}

func ParseFileSystem(name string, text []byte, fs http.FileSystem) (string, error) {
	outTpl, err := NewFileSystem(name, fs).Parse(text)
	if err != nil {
		return "", err
	}
	b := new(bytes.Buffer)
	outTpl.WriteIn(b)
	return b.String(), nil
}

// ParseFile parse the jade template file in given filename
func ParseFile(filename string) (string, error) {
	bs, err := ReadFunc(filename, nil)
	if err != nil {
		return "", err
	}
	return Parse(filepath.Base(filename), bs)
}

func (t *tree) WriteIn(b io.Writer) {
	t.Root.WriteIn(b)
}

func ReadFunc(filename string, fs http.FileSystem) ([]byte, error) {
	if fs == nil {
		return ioutil.ReadFile(filename)
	}

	file, err := fs.Open(filename[1:])
	if err != nil {
		return nil, err
	}

	info, err := file.Stat()
	if err != nil {
		return nil, err
	}

	buf := make([]byte, info.Size())
	_, err = file.Read(buf)
	if err != nil {
		return nil, err
	}

	return buf, nil
}
