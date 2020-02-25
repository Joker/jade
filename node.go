// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package jade

import (
	"bytes"
	"io"
)

// var textFormat = "%s" // Changed to "%q" in tests for better error messages.

// A Node is an element in the parse tree. The interface is trivial.
// The interface contains an unexported method so that only
// types local to this package can satisfy it.
type Node interface {
	Type() NodeType
	String() string
	WriteIn(io.Writer)
	// Copy does a deep copy of the Node and all its components.
	// To avoid type assertions, some XxxNodes also have specialized
	// CopyXxx methods that return *XxxNode.
	Copy() Node
	Position() Pos // byte position of start of node in full original input string
	// tree returns the containing *Tree.
	// It is unexported so all implementations of Node are in this package.
	tree() *Tree
}

// Pos represents a byte position in the original input text from which
// this template was parsed.
type Pos int

func (p Pos) Position() Pos {
	return p
}

// Nodes.

// listNode holds a sequence of nodes.
type listNode struct {
	NodeType
	Pos
	tr    *Tree
	Nodes []Node // The element nodes in lexical order.
}

func (t *Tree) newList(pos Pos) *listNode {
	return &listNode{tr: t, NodeType: NodeList, Pos: pos}
}

func (l *listNode) append(n Node) {
	l.Nodes = append(l.Nodes, n)
}

func (l *listNode) tree() *Tree {
	return l.tr
}

func (l *listNode) String() string {
	b := new(bytes.Buffer)
	l.WriteIn(b)
	return b.String()
}
func (l *listNode) WriteIn(b io.Writer) {
	for _, n := range l.Nodes {
		n.WriteIn(b)
	}
}

func (l *listNode) CopyList() *listNode {
	if l == nil {
		return l
	}
	n := l.tr.newList(l.Pos)
	for _, elem := range l.Nodes {
		n.append(elem.Copy())
	}
	return n
}

func (l *listNode) Copy() Node {
	return l.CopyList()
}
