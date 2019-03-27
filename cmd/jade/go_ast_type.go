package main

import (
	"go/ast"
	"go/importer"
	"go/types"
	"strconv"
	"strings"
)

func (a *goAST) checkType() {
	var (
		info = types.Info{
			Types: make(map[ast.Expr]types.TypeAndValue),
			Defs:  make(map[*ast.Ident]types.Object),
			Uses:  make(map[*ast.Ident]types.Object),
		}
		conf = types.Config{
			Importer: importer.Default(),
			// DisableUnusedImportCheck: true,
			Error: func(err error) {
				// fmt.Println(err)
			},
		}
	)
	conf.Check("check", a.fset, []*ast.File{a.node}, &info)

	ast.Inspect(a.node, func(n ast.Node) bool {
		if n != nil {
			switch x := n.(type) {
			case *ast.CaseClause:
				rewrite(x.Body, &info)
			case *ast.BlockStmt:
				rewrite(x.List, &info)
				// case *ast.CallExpr:
				// return false
			}
		}
		return true
	})
}

//

func rewrite(in []ast.Stmt, info *types.Info) {
	for k, ex := range in {
		if d, ok := ex.(*ast.DeclStmt); ok {
			if s, ok := d.Decl.(*ast.GenDecl); ok && len(s.Specs) == 1 {
				if v, ok := s.Specs[0].(*ast.ValueSpec); ok && len(v.Names) == 1 {
					var escape bool

					switch {
					case strings.HasPrefix(v.Names[0].Name, "esc"):
						escape = true
					case strings.HasPrefix(v.Names[0].Name, "unesc"):
						escape = false
					default:
						continue
					}

					switch vt := info.TypeOf(v.Values[0]).(type) {
					case *types.Basic:
						if stdlib {
							switch vt.Name() {
							case "string":
								in[k] = stdlibFuncCall(escape, "", "", arg(v.Values[0]))
							case "int", "int8", "int16", "int32":
								in[k] = stdlibFuncCall(escape, "strconv", "FormatInt", arg(funcCall("", "int64", arg(v.Values[0])), a("10")))
							case "int64":
								in[k] = stdlibFuncCall(escape, "strconv", "FormatInt", arg(v.Values[0], a("10")))
							case "uint", "uint8", "uint16", "uint32":
								in[k] = stdlibFuncCall(escape, "strconv", "FormatUint", arg(funcCall("", "uint64", arg(v.Values[0])), a("10")))
							case "uint64":
								in[k] = stdlibFuncCall(escape, "strconv", "FormatUint", arg(v.Values[0], a("10")))
							case "bool":
								in[k] = stdlibFuncCall(escape, "strconv", "FormatBool", arg(v.Values[0]))
							case "float64":
								in[k] = stdlibFuncCall(escape, "strconv", "FormatFloat", arg(v.Values[0], a("'f'"), &ast.UnaryExpr{Op: 13, X: a("1")}, a("64"))) // &ast.UnaryExpr{Op: 13, X: a("1")} =>= -1
							default:
								in[k] = stdlibFuncCall(escape, "fmt", "Sprintf", arg(a(`"%v"`), v.Values[0]))
							}
						} else {
							switch vt.Name() {
							case "string":
								if escape {
									in[k] = &ast.ExprStmt{X: funcCall(lib_name, "WriteEscString", arg(v.Values[0], a("buffer")))}
								} else {
									in[k] = &ast.ExprStmt{X: funcCall("buffer", "WriteString", arg(v.Values[0]))}
								}
							case "int", "int8", "int16", "int32":
								in[k] = &ast.ExprStmt{X: funcCall(lib_name, "WriteInt", arg(funcCall("", "int64", arg(v.Values[0])), a("buffer")))}
							case "int64":
								in[k] = &ast.ExprStmt{X: funcCall(lib_name, "WriteInt", arg(v.Values[0], a("buffer")))}
							case "uint", "uint8", "uint16", "uint32":
								in[k] = &ast.ExprStmt{X: funcCall(lib_name, "WriteUint", arg(funcCall("", "uint64", arg(v.Values[0])), a("buffer")))}
							case "uint64":
								in[k] = &ast.ExprStmt{X: funcCall(lib_name, "WriteUint", arg(v.Values[0], a("buffer")))}
							case "bool":
								in[k] = &ast.ExprStmt{X: funcCall(lib_name, "WriteBool", arg(v.Values[0], a("buffer")))}
							case "float64":
								in[k] = &ast.ExprStmt{X: funcCall(lib_name, "WriteFloat", arg(v.Values[0], a("buffer")))}
							default:
								in[k] = &ast.ExprStmt{X: funcCall(lib_name, "WriteAll", arg(v.Values[0], a(strconv.FormatBool(escape)), a("buffer")))}
							}
						}
					default:
						if stdlib {
							in[k] = stdlibFuncCall(escape, "fmt", "Sprintf", arg(a(`"%v"`), v.Values[0]))
						} else {
							in[k] = &ast.ExprStmt{X: funcCall(lib_name, "WriteAll", arg(v.Values[0], a(strconv.FormatBool(escape)), a("buffer")))}
						}
					}
				}
			}
		}
	}
}

func stdlibFuncCall(esc bool, x, sel string, exps []ast.Expr) *ast.ExprStmt {
	arg := exps
	if sel != "" {
		arg = []ast.Expr{funcCall(x, sel, exps)}
	}
	if esc {
		return &ast.ExprStmt{X: funcCall("buffer", "WriteString", []ast.Expr{funcCall("html", "EscapeString", arg)})}
	} else {
		return &ast.ExprStmt{X: funcCall("buffer", "WriteString", arg)}
	}
}

//

func funcCall(packName, funcName string, exps []ast.Expr) *ast.CallExpr {
	if packName == "" {
		return &ast.CallExpr{
			Fun: &ast.Ident{
				Name: funcName,
			},
			Args: exps,
		}
	}
	return &ast.CallExpr{
		Fun: &ast.SelectorExpr{
			X:   &ast.Ident{Name: packName},
			Sel: &ast.Ident{Name: funcName},
		},
		Args: exps,
	}
}

func arg(i ...ast.Expr) []ast.Expr { return i }
func a(i string) *ast.BasicLit     { return &ast.BasicLit{Value: i} }
