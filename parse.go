// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package jade

import (
	"fmt"
	"runtime"
	"strings"
)

// Tree is the representation of a single parsed template.
type Tree struct {
	Name      string    // name of the template represented by the tree.
	ParseName string    // name of the top-level template during parsing, for error messages.
	Root      *listNode // top-level root of the tree.
	text      string    // text parsed to create the template (or its parent)
	// Parsing only; cleared after parse.
	funcs     []map[string]interface{}
	lex       *lexer
	token     [3]item // three-token lookahead for parser.
	peekCount int
	vars      []string // variables defined at the moment.
}

// Copy returns a copy of the Tree. Any parsing state is discarded.
func (t *Tree) Copy() *Tree {
	if t == nil {
		return nil
	}
	return &Tree{
		Name:      t.Name,
		ParseName: t.ParseName,
		Root:      t.Root.CopyList(),
		text:      t.text,
	}
}

// Parse returns a map from template name to parse.Tree, created by parsing the
// templates described in the argument string. The top-level template will be
// given the specified name. If an error is encountered, parsing stops and an
// empty map is returned with the error.
/*
func Parse(name, text, leftDelim, rightDelim string, funcs ...map[string]interface{}) (treeSet map[string]*Tree, err error) {
	treeSet = make(map[string]*Tree)
	t := New(name)
	t.text = text
	_, err = t.Parse(text, leftDelim, rightDelim, treeSet, funcs...)
	return
}
// */

// next returns the next token.
func (t *Tree) next() item {
	if t.peekCount > 0 {
		t.peekCount--
	} else {
		t.token[0] = t.lex.nextItem()
	}
	return t.token[t.peekCount]
}

// backup backs the input stream up one token.
func (t *Tree) backup() {
	t.peekCount++
}

// backup2 backs the input stream up two tokens.
// The zeroth token is already there.
func (t *Tree) backup2(t1 item) {
	t.token[1] = t1
	t.peekCount = 2
}

// backup3 backs the input stream up three tokens
// The zeroth token is already there.
func (t *Tree) backup3(t2, t1 item) { // Reverse order: we're pushing back.
	t.token[1] = t1
	t.token[2] = t2
	t.peekCount = 3
}

// peek returns but does not consume the next token.
func (t *Tree) peek() item {
	if t.peekCount > 0 {
		return t.token[t.peekCount-1]
	}
	t.peekCount = 1
	t.token[0] = t.lex.nextItem()
	return t.token[0]
}

/*
// nextNonSpace returns the next non-space token.
func (t *Tree) nextNonSpace() (token item) {
	for {
		token = t.next()
		if token.typ != itemSpace {
			break
		}
	}
	return token
}

// peekNonSpace returns but does not consume the next non-space token.
func (t *Tree) peekNonSpace() (token item) {
	for {
		token = t.next()
		if token.typ != itemSpace {
			break
		}
	}
	t.backup()
	return token
}
// */

// Parsing.

// New allocates a new parse tree with the given name.
func New(name string, funcs ...map[string]interface{}) *Tree {
	return &Tree{
		Name:  name,
		funcs: funcs,
	}
}

// ErrorContext returns a textual representation of the location of the node in the input text.
// The receiver is only used when the node does not have a pointer to the tree inside,
// which can occur in old code.
func (t *Tree) ErrorContext(n Node) (location, context string) {
	pos := int(n.Position())
	tree := n.tree()
	if tree == nil {
		tree = t
	}
	text := tree.text[:pos]
	byteNum := strings.LastIndex(text, "\n")
	if byteNum == -1 {
		byteNum = pos // On first line.
	} else {
		byteNum++ // After the newline.
		byteNum = pos - byteNum
	}
	lineNum := 1 + strings.Count(text, "\n")
	context = n.String()
	if len(context) > 20 {
		context = fmt.Sprintf("%.20s...", context)
	}
	return fmt.Sprintf("%s:%d:%d", tree.ParseName, lineNum, byteNum), context
}

// errorf formats the error and terminates processing.
func (t *Tree) errorf(format string, args ...interface{}) {
	t.Root = nil
	format = fmt.Sprintf("template: %s:%d: %s", t.ParseName, t.lex.lineNumber(), format)
	panic(fmt.Errorf(format, args...))
}

// error terminates processing.
func (t *Tree) error(err error) {
	t.errorf("%s", err)
}

/*
// expect consumes the next token and guarantees it has the required type.
func (t *Tree) expect(expected itemType, context string) item {
	token := t.nextNonSpace()
	if token.typ != expected {
		t.unexpected(token, context)
	}
	return token
}

// expectOneOf consumes the next token and guarantees it has one of the required types.
func (t *Tree) expectOneOf(expected1, expected2 itemType, context string) item {
	token := t.nextNonSpace()
	if token.typ != expected1 && token.typ != expected2 {
		t.unexpected(token, context)
	}
	return token
}
// */

// unexpected complains about the token and terminates processing.
func (t *Tree) unexpected(token item, context string) {
	t.errorf("unexpected %s in %s", token, context)
}

// recover is the handler that turns panics into returns from the top level of Parse.
func (t *Tree) recover(errp *error) {
	e := recover()
	if e != nil {
		if _, ok := e.(runtime.Error); ok {
			panic(e)
		}
		if t != nil {
			t.stopParse()
		}
		*errp = e.(error)
	}
	return
}

// startParse initializes the parser, using the lexer.
func (t *Tree) startParse(funcs []map[string]interface{}, lex *lexer) {
	t.Root = nil
	t.lex = lex
	t.vars = []string{"$"}
	t.funcs = funcs
}

// stopParse terminates parsing.
func (t *Tree) stopParse() {
	t.lex = nil
	t.vars = nil
	t.funcs = nil
}

// Parse parses the template definition string to construct a representation of
// the template for execution. If either action delimiter string is empty, the
// default ("{{" or "}}") is used. Embedded template definitions are added to
// the treeSet map.
func (t *Tree) Parse(text, leftDelim, rightDelim string, treeSet map[string]*Tree, funcs ...map[string]interface{}) (tree *Tree, err error) {
	defer t.recover(&err)
	t.ParseName = t.Name
	t.startParse(funcs, lex(t.Name, text, leftDelim, rightDelim))
	t.text = text
	t.parse(treeSet)
	t.add(treeSet)
	t.stopParse()
	return t, nil
}

// add adds tree to the treeSet.
func (t *Tree) add(treeSet map[string]*Tree) {
	tree := treeSet[t.Name]
	if tree == nil || IsEmptyTree(tree.Root) {
		treeSet[t.Name] = t
		return
	}
	if !IsEmptyTree(t.Root) {
		t.errorf("template: multiple definition of template %q", t.Name)
	}
}

// IsEmptyTree reports whether this tree (node) is empty of everything but space.
func IsEmptyTree(n Node) bool {
	switch n := n.(type) {
	case nil:
		return true
	// case *ActionNode:
	case *listNode:
		for _, node := range n.Nodes {
			if !IsEmptyTree(node) {
				return false
			}
		}
		return true
	// case *TextNode:
	// 	return len(bytes.TrimSpace(n.Text)) == 0
	default:
		panic("unknown node: " + n.String())
	}
	return false
}
