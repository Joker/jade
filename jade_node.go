package jade

import (
	"bytes"
	"fmt"
)



// NodeType identifies the type of a parse tree node.
type NodeType int
const (
	NodeList    NodeType = iota
	NodeTag
)



type TagNode struct {
	NodeType
	Pos
	tr    *Tree
	Nodes []Node

	Tag     string
	Indent  int
	Nesting int
}

func (t *Tree) newTag(pos Pos, tag string) *TagNode {
	return &TagNode{tr: t, NodeType: NodeTag, Pos: pos, Tag: tag}
}

func (l *TagNode) append(n Node) {
	l.Nodes = append(l.Nodes, n)
}

func (l *TagNode) tree() *Tree {
	return l.tr
}

func (l *TagNode) String() string {
	b := new(bytes.Buffer)

	fmt.Fprint(b, fmt.Sprintf("<%s>\n", l.Tag))

	for _, n := range l.Nodes {
		if n.Type() == NodeTag { fmt.Fprint(b, n) }
	}

	fmt.Fprintf(b, "</%s>\n", l.Tag)
	return b.String()
}

func (l *TagNode) CopyTag() *TagNode {
	if l == nil {
		return l
	}
	n := l.tr.newTag(l.Pos, string(l.Tag))
	for _, elem := range l.Nodes {
		n.append(elem.Copy())
	}
	return n
}

func (l *TagNode) Copy() Node {
	return l.CopyTag()
}
