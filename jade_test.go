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
		fmt.Println("_________" + name + ".jade")

		dat, err := ioutil.ReadFile("testdata/" + name + ".jade")
		if err != nil {
			fmt.Printf("--- FAIL: ReadFile error: %v\n\n", err)
			t.Fail()
			continue
		}

		tpl, err := Parse(name+".jade", string(dat))
		if err != nil {
			fmt.Printf("--- FAIL: Parse error: %v\n\n", err)
			t.Fail()
			continue
		}
		tmpl := bufio.NewScanner(strings.NewReader(tpl))
		tmpl.Split(bufio.ScanLines)

		inFile, err := os.Open("testdata/" + name + ".html")
		if err != nil {
			fmt.Printf("--- FAIL: OpenFile error: %v\n\n", err)
			t.Fail()
			continue
		}
		file := bufio.NewScanner(inFile)
		file.Split(bufio.ScanLines)

		nilerr := 0
		line := 0
		for tmpl.Scan() {
			file.Scan()

			a := tmpl.Text()
			b := file.Text()
			line += 1

			if strings.Compare(a, b) != 0 && nilerr < 4 {
				fmt.Printf("%s\n%s\n%d^___________________________\n", a, b, line)
				nilerr += 1
				t.Fail()
			}
		}
		inFile.Close()

		if nilerr == 0 {
			fmt.Println("    PASS\n")
		} else {
			fmt.Println("--- FAIL\n")
		}
	}
	// ioutil.WriteFile("testdata/"+name+".html", []byte(tpl), 0644)
}
