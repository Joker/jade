package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/Joker/jade"
	"golang.org/x/tools/imports"
)

var (
	dict     = map[string]string{}
	lib_name = ""
	outdir   string
	pkg_name string
	stdlib   bool
	stdbuf   bool
	inline   bool
	format   bool
)

func use() {
	fmt.Printf("Usage: %s [OPTION]... [FILE]... \n", os.Args[0])
	flag.PrintDefaults()
}
func init() {
	flag.StringVar(&outdir, "d", "./", `directory for generated .go files`)
	flag.StringVar(&pkg_name, "pkg", "jade", `package name for generated files`)
	flag.BoolVar(&format, "fmt", false, `HTML pretty print output for generated functions`)
	flag.BoolVar(&inline, "inline", false, `inline HTML in generated functions`)
	flag.BoolVar(&stdlib, "stdlib", false, `use stdlib functions`)
	flag.BoolVar(&stdbuf, "stdbuf", false, `use bytes.Buffer  [default bytebufferpool.ByteBuffer]`)
}

//

type goAST struct {
	node *ast.File
	fset *token.FileSet
}

func (a *goAST) bytes(bb *bytes.Buffer) []byte {
	printer.Fprint(bb, a.fset, a.node)
	return bb.Bytes()
}

func parseGoSrc(fileName string, GoSrc interface{}) (out goAST, err error) {
	out.fset = token.NewFileSet()
	out.node, err = parser.ParseFile(out.fset, fileName, GoSrc, parser.ParseComments)
	return
}

func goImports(absPath string, src []byte) []byte {
	fmtOut, err := imports.Process(absPath, src, &imports.Options{TabWidth: 4, TabIndent: true, Comments: true, Fragment: true})
	if err != nil {
		log.Fatalln("goImports(): ", err)
	}

	return fmtOut
}

//

func genFile(path, outdir, pkg_name string) {
	log.Printf("file: %q\n", path)

	var (
		dir       = filepath.Dir(path)
		fname     = filepath.Base(path)
		outPath   = outdir + "/" + fname + ".go"
		rx, _     = regexp.Compile("[^a-zA-Z0-9]+")
		constName = rx.ReplaceAllString(fname[:len(fname)-4], "")
	)
	if wd, err := os.Getwd(); err == nil && wd != dir {
		os.Chdir(dir)
	}

	fl, err := ioutil.ReadFile(fname)
	if err != nil {
		log.Fatalln("jade: ReadFile(): ", err)
	}

	//

	jst, err := jade.New(path).Parse(fl)
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

	//

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

	err = ioutil.WriteFile(outPath, fmtOut, 0644)
	if err != nil {
		log.Fatalln("jade: WriteFile(): ", err)
	}
}

func genDir(dir, outdir, pkg_name string) {
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("prevent panic by handling failure accessing a path %q: %v\n", dir, err)
		}

		if ext := filepath.Ext(info.Name()); ext == ".jade" || ext == ".pug" {
			genFile(path, outdir, pkg_name)
		}
		return nil
	})
	if err != nil {
		log.Fatalln(err)
	}
}

//

func main() {
	flag.Usage = use
	flag.Parse()
	if len(flag.Args()) == 0 {
		use()
		return
	}

	jade.Config(golang)

	if _, err := os.Stat(outdir); os.IsNotExist(err) {
		os.MkdirAll(outdir, 0755)
	}
	outdir, _ = filepath.Abs(outdir)

	for _, jadePath := range flag.Args() {

		stat, err := os.Stat(jadePath)
		if err != nil {
			log.Fatalln(err)
		}

		absPath, _ := filepath.Abs(jadePath)
		if stat.IsDir() {
			genDir(absPath, outdir, pkg_name)
		} else {
			genFile(absPath, outdir, pkg_name)
		}
		if !stdlib {
			makeJfile(stdbuf)
		}
	}
}
