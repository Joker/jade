package jade

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Joker/hpp"
)

func examination(test func(dat string) (string, error), ext, path string, t *testing.T) {

	files, _ := ioutil.ReadDir(path)

	var name, fext string
	for _, file := range files {
		name = file.Name()
		fext = filepath.Ext(name)

		if fext != ".jade" && fext != ".pug" {
			continue
		}

		dat, err := ioutil.ReadFile(path + name)
		if err != nil {
			fmt.Println("_________" + name)
			fmt.Printf("--- FAIL: ReadFile error: %v\n\n", err)
			t.Fail()
			continue
		}

		tpl, err := test(string(dat))
		if err != nil {
			fmt.Println("_________" + name)
			fmt.Printf("--- FAIL: test run() error: \n%s\n\n", err)
			t.Fail()
			continue
		}

		tmpl := bufio.NewScanner(strings.NewReader(tpl))
		tmpl.Split(bufio.ScanLines)

		inFile, err := os.Open(path + strings.TrimSuffix(name, fext) + ext)
		if err != nil {

			// make files
			ioutil.WriteFile(path+strings.TrimSuffix(name, fext)+ext, []byte(tpl), 0644)

			fmt.Println("```", tpl, "\n\n```")
			continue
		}

		html := bufio.NewScanner(inFile)
		html.Split(bufio.ScanLines)

		nilerr, line := 0, 0

		for tmpl.Scan() {
			html.Scan()

			a := tmpl.Text()
			b := html.Text()
			line += 1

			if strings.Compare(a, b) != 0 && nilerr < 4 {
				if nilerr == 0 {
					fmt.Println("_________" + name)
				}
				fmt.Printf("%s\n%s\n%d^___________________________\n", a, b, line)
				nilerr += 1
				t.Fail()
			}
		}
		inFile.Close()

		if nilerr != 0 {
			fmt.Print("--- FAIL\n\n")
		}
		// } else { fmt.Print("    PASS\n\n") }
	}
}

func lexerTest(dat string) (string, error) {
	var buf bytes.Buffer

	l := lex("test", dat)
	for i := range l.items {
		switch {
		case i.typ == itemError:
			buf.WriteString("\n\n\nError:\n\t")
			buf.WriteString(fmt.Sprintf("%s  line: %d", i.val, i.line))
			buf.WriteString("\n\n\n")
			return "", fmt.Errorf("%s", buf.String())
		case i.typ == itemEOF:
			buf.WriteString("\nEOF")
		case i.typ == itemEndL:
			buf.WriteByte('\n')
		case i.typ == itemEmptyLine:
			buf.WriteString(i.val)
		case i.typ == itemIdent:
			buf.WriteString(i.val)
		default:
			buf.WriteString(fmt.Sprintf("[%d  %s   \"%s\"]\t", i.depth, i.typ, i.val))
		}
	}
	return buf.String(), nil
}

func xTestJadeLex(t *testing.T) {
	examination(lexerTest, ".lex", "./testdata/", t)
}

//

func parserTest(text string) (string, error) {
	outTpl, err := New("test").Parse(text)
	if err != nil {
		return "", err
	}
	b := new(bytes.Buffer)
	outTpl.WriteIn(b)
	return string(hpp.Print(b)), nil
}

func TestJadeParse(t *testing.T) {
	examination(parserTest, ".html", "./testdata/", t)
}
