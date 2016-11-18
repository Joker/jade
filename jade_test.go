package jade

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

var fname = []string{
	"attributes",
	"code",
	"dynamicscript",
	"each",
	"extend",
	"extend-layout",
	"form",
	"includes",
	"layout",
	"mixins",
	"pet",
	"rss",
	"text",
	"whitespace",
}

func gen(t *testing.T) {
	for _,name := range fname {
		dat, err := ioutil.ReadFile("testdata/"+name+".jade")
		if err != nil {
			t.Error("ReadFile error: %v", err)
			return
		}

		tpl, err := Parse("tpl_"+name, string(dat))
		if err != nil {
			t.Error("Parse error: %v", err)
			return
		}

	    ioutil.WriteFile("testdata/"+name+".tmpl", []byte(tpl), 0644)
	}
}


func TestAttributes(t *testing.T) {
	dat, err := ioutil.ReadFile("testdata/attributes.jade")
	if err != nil {
		t.Error("ReadFile error: %v", err)
		return
	}

	tpl, err := Parse("tpl", string(dat))
	if err != nil {
		t.Error("Parse error: %v", err)
		return
	}

	tmpl := bufio.NewScanner(strings.NewReader(tpl))
	tmpl.Split(bufio.ScanLines)

	inFile, _ := os.Open("testdata/attributes.tmpl")
	defer inFile.Close()
	file := bufio.NewScanner(inFile)
	file.Split(bufio.ScanLines)

	for tmpl.Scan() {
		file.Scan()

		a := tmpl.Text()
		b := file.Text()

		if strings.Compare(a, b) != 0 {
			fmt.Printf("%s\n%s\n________________________\n", a, b)
			t.Fail()
		}
	}
}
