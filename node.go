// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package jade

import (
	"bytes"
	"fmt"
)

var textFormat = "%s" // Changed to "%q" in tests for better error messages.

// A Node is an element in the parse tree. The interface is trivial.
// The interface contains an unexported method so that only
// types local to this package can satisfy it.
type Node interface {
	Type() NodeType
	position() Pos // byte position of start of node in full original input string
	String() string

	// Copy does a deep copy of the Node and all its components.
	// To avoid type assertions, some XxxNodes also have specialized
	// CopyXxx methods that return *XxxNode.
	Copy() Node

	// tree returns the containing *Tree.
	// It is unexported so all implementations of Node are in this package.
	tree() *Tree
	tp() itemType
}

// Type returns itself and provides an easy default implementation
// for embedding in a Node. Embedded in all non-trivial Nodes.
func (t NodeType) Type() NodeType {
	return t
}

// Pos represents a byte position in the original input text from which
// this template was parsed.
type Pos int

func (p Pos) position() Pos {
	return p
}








// listNode holds a sequence of nodes.
type listNode struct {
	NodeType
	Pos
	tr    *Tree
	Nodes []Node // The element nodes in lexical order.
}

func (t *Tree) newList(pos Pos) *listNode {
	return &listNode{tr: t, NodeType: nodeList, Pos: pos}
}

func (l *listNode) append(n Node) {
	l.Nodes = append(l.Nodes, n)
}

func (l *listNode) tree() *Tree {
	return l.tr
}
func (l *listNode) tp() itemType {
	return 0
}

func (l *listNode) String() string {
	b := new(bytes.Buffer)
	for _, n := range l.Nodes {
		fmt.Fprint(b, n)
	}
	return b.String()
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
