package main

import (
	"go/ast"
)

func (a *goAST) checkUnresolvedBlock() {
	var dsSl = []*ast.DeclStmt{}
	var ds *ast.DeclStmt

	ast.Inspect(a.node, func(n ast.Node) bool {
		if n != nil {
			switch xds := n.(type) {
			case *ast.DeclStmt:
				if gd, ok := xds.Decl.(*ast.GenDecl); ok {
					ast.Inspect(gd, func(n ast.Node) bool {
						if n != nil {
							switch x := n.(type) {
							case *ast.Ident:
								if x.Name == "block" {
									dsSl = append(dsSl, xds)
									ds = xds
								}
							case *ast.CallExpr:
								return false
							}
						}
						return true
					})
				}
			case *ast.CallExpr:
				if len(xds.Args) == 1 {
					if i, ok := xds.Args[0].(*ast.Ident); ok {
						if i.Name == "block" {
							if dsSl[len(dsSl)-1] == ds {
								dsSl = dsSl[:len(dsSl)-1]
							}
						}
					}
				}
			}
		}
		return true
	})
	ast.Inspect(a.node, func(n ast.Node) bool {
		if n != nil {
			switch x := n.(type) {
			case *ast.CallExpr:
				return false
			case *ast.BlockStmt:
				for i, v := range x.List {
					if ds, ok := v.(*ast.DeclStmt); ok {
						for _, v := range dsSl {
							if ds == v {
								x.List = append(x.List[:i], x.List[i+1:]...)
							}
						}
					}
				}
			}
		}
		return true
	})
}
