package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Joker/jade"
)

var wdir string

func init() {
	os.Chdir("../..")
	wdir, _ = os.Getwd()
}

func examination(test func(dat []byte) ([]byte, error), ext, path string, t *testing.T) {
	os.Chdir(path)
	files, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Printf("--- FAIL: ReadDir error: %v\n\n", err)
		t.Fail()
	}

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

		tpl, err := test(dat)
		if err != nil {
			fmt.Println("_________" + name)
			fmt.Printf("--- FAIL: test run() error: \n%s\n\n", err)
			t.Fail()
			continue
		}

		tmpl := bufio.NewScanner(bytes.NewReader(tpl))
		tmpl.Split(bufio.ScanLines)

		inFile, err := os.Open(path + strings.TrimSuffix(name, fext) + ext)
		if err != nil {

			// make files
			ioutil.WriteFile(path+strings.TrimSuffix(name, fext)+ext, []byte(tpl), 0644)

			fmt.Println("```", string(tpl), "\n\n```")
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
					fmt.Println("_________" + name + "\n")
				}
				fmt.Printf("%s\n%s\n%d^___________________________\n", a, b, line)
				nilerr += 1
				t.Fail()
			}
		}
		inFile.Close()

		if nilerr != 0 {
			fmt.Print("--- FAIL\n\n\n\n")
		}
	}
}

func astTest(text []byte) ([]byte, error) {
	jade.ConfigOtputGo()

	constName := "test"
	outPath := "test"
	inline = true

	//

	jst, err := jade.New("path").Parse(text)
	if err != nil {
		log.Fatalln("jade: jade.New(path).Parse(): ", err)
	}

	var (
		bb  = new(bytes.Buffer)
		tpl = newLayout(constName)
	)
	tpl.writeBefore(bb)
	jst.WriteIn(bb)
	tpl.writeAfter(bb)

	gst, err := parseGoSrc(outPath, bb)
	if err != nil {
		log.Fatalln("jade: parseGoSrc(): ", err)
	}

	gst.collapseWriteString(inline, constName)
	gst.checkType()
	gst.checkUnresolvedBlock()

	bb.Reset()
	fmtOut := goImports(outPath, gst.bytes(bb))

	//

	return fmtOut, nil
}

func TestGoASToptimize(t *testing.T) {
	examination(astTest, ".go", wdir+"/testdata/v2/", t)
	examination(astTest, ".go", wdir+"/testdata/v2/includes/", t)
	examination(astTest, ".go", wdir+"/testdata/v2/inheritance/", t)
}
