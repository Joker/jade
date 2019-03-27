package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"strings"
)

func (a *goAST) collapseWriteString(strInline bool, constName string) {
	// ast.Print(a.fset, a.node)

	var wrSt = writeStrings{inline: strInline, constName: constName}
	ast.Inspect(a.node, func(n ast.Node) bool {
		if n != nil {
			switch x := n.(type) {
			case *ast.CaseClause:
				wrSt.push(x.Body)
			case *ast.BlockStmt:
				wrSt.push(x.List)
				// case *ast.CallExpr:
				// return false
			}
		}
		return true
	})

	wrSt.fillConstNode(a.node.Decls)
}

//

type writeStrings struct {
	constSlice []struct {
		str  string
		name string
	}
	constName  string
	constCount int
	inline     bool
}

func (ws *writeStrings) push(in []ast.Stmt) {
	type wStr struct {
		s string
		a *ast.CallExpr
		z *ast.Stmt
	}
	List := [][]wStr{}
	subList := make([]wStr, 0, 4)

	for k, ex := range in {
		if es, ok := ex.(*ast.ExprStmt); ok {
			if ce, ok := es.X.(*ast.CallExpr); ok {
				if fun, ok := ce.Fun.(*ast.SelectorExpr); ok && len(ce.Args) == 1 && fun.Sel.Name == "WriteString" {
					if arg, ok := ce.Args[0].(*ast.BasicLit); ok {
						subList = append(subList, wStr{s: strings.Trim(arg.Value, "`"), a: ce, z: &in[k]})
						continue
					}
				}
			}
		}
		if len(subList) > 0 {
			List = append(List, subList)
			subList = make([]wStr, 0, 4)
		}
	}
	if len(subList) > 0 {
		List = append(List, subList)
	}

	//

	var st = new(bytes.Buffer)
	for _, block := range List {

		st.WriteString(block[0].s)
		for i := 1; i < len(block); i++ {
			st.WriteString(block[i].s)
			*block[i].z = new(ast.EmptyStmt) // remove a node
		}

		if ws.inline {
			block[0].a.Args[0].(*ast.BasicLit).Value = "`" + st.String() + "`"
		} else {
			str := st.String()
			if name, ok := dict[str]; ok {
				block[0].a.Args = []ast.Expr{&ast.Ident{Name: name}}
			} else {
				newName := fmt.Sprintf("%s__%d", ws.constName, ws.constCount)
				block[0].a.Args = []ast.Expr{&ast.Ident{Name: newName}}
				dict[str] = newName
				ws.constSlice = append(ws.constSlice, struct {
					str  string
					name string
				}{
					str,
					newName,
				})
			}
		}

		st.Reset()
		ws.constCount += 1
	}
}

func (ws *writeStrings) fillConstNode(decl []ast.Decl) {
	if constNode, ok := decl[1].(*ast.GenDecl); ok && !ws.inline {
		for _, v := range ws.constSlice {
			constNode.Specs = append(constNode.Specs, &ast.ValueSpec{
				Names: []*ast.Ident{
					&ast.Ident{Name: v.name},
				},
				Values: []ast.Expr{
					&ast.BasicLit{Kind: 9, Value: "`" + v.str + "`"}, // 9 => string
				},
			})
		}
	}
}
