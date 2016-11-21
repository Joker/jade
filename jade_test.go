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

func TestJadeExamples(t *testing.T) {
	for _, name := range fname {
		fmt.Println("========= testing: " + name + ".jade")

		dat, err := ioutil.ReadFile("testdata/" + name + ".jade")
		if err != nil {
			fmt.Printf("--- FAIL: ReadFile error: %v\n", err)
			t.Fail()
			continue
		}

		tpl, err := Parse("tpl"+name, string(dat))
		if err != nil {
			fmt.Printf("--- FAIL: Parse error: %v\n", err)
			t.Fail()
			continue
		}
		tmpl := bufio.NewScanner(strings.NewReader(tpl))
		tmpl.Split(bufio.ScanLines)

		inFile, err := os.Open("testdata/" + name + ".html")
		if err != nil {
			fmt.Printf("--- FAIL: OpenFile error: %v\n", err)
			t.Fail()
			continue
		}
		file := bufio.NewScanner(inFile)
		file.Split(bufio.ScanLines)

		nilerr := true
		for tmpl.Scan() {
			file.Scan()

			a := tmpl.Text()
			b := file.Text()

			if strings.Compare(a, b) != 0 {
				fmt.Printf("%s\n%s\n___________________________\n", a, b)
				nilerr = false
				t.Fail()
			}
		}
		inFile.Close()

		if nilerr {
			fmt.Println("--- PASS")
		} else {
			fmt.Println("--- FAIL")
		}
	}
	// ioutil.WriteFile("testdata/"+name+".html", []byte(tpl), 0644)
}
