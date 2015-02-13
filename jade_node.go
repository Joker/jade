package jade

import (
	"bytes"
	"fmt"
)



// NodeType identifies the type of a parse tree node.
type NodeType int
const (
	NodeList    NodeType = iota
	NodeText
	NodeTag
	NodeDoctype
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

func (tg *TagNode) append(n Node) {
	tg.Nodes = append(tg.Nodes, n)
}
func (tg *TagNode) tree() *Tree {
	return tg.tr
}

func (tg *TagNode) String() string {
	b := new(bytes.Buffer)

	fmt.Fprint(b, fmt.Sprintf("<%s>\n", tg.Tag))

	for _, n := range tg.Nodes {
		if n.Type() == NodeTag || n.Type() == NodeText { fmt.Fprint(b, n) }
	}

	fmt.Fprintf(b, "</%s>\n", tg.Tag)
	return b.String()
}

func (tg *TagNode) CopyTag() *TagNode {
	if tg == nil {
		return tg
	}
	n := tg.tr.newTag(tg.Pos, string(tg.Tag))
	for _, elem := range tg.Nodes {
		n.append(elem.Copy())
	}
	return n
}
func (tg *TagNode) Copy() Node {
	return tg.CopyTag()
}



var doctype = map[string]string {
	"xml" 			: `<?xml version="1.0" encoding="utf-8" ?>`,
	"html" 			: `<!DOCTYPE html>`,
	"5" 			: `<!DOCTYPE html>`,
	"1.1" 			: `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.1//EN" "http://www.w3.org/TR/xhtml11/DTD/xhtml11.dtd">`,
	"xhtml" 		: `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.1//EN" "http://www.w3.org/TR/xhtml11/DTD/xhtml11.dtd">`,
	"basic" 		: `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML Basic 1.1//EN" "http://www.w3.org/TR/xhtml-basic/xhtml-basic11.dtd">`,
	"mobile" 		: `<!DOCTYPE html PUBLIC "-//WAPFORUM//DTD XHTML Mobile 1.2//EN" "http://www.openmobilealliance.org/tech/DTD/xhtml-mobile12.dtd">`,
	"strict" 		: `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Strict//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd">`,
	"frameset" 		: `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Frameset//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-frameset.dtd">`,
	"transitional" 	: `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">`,
	"4" 			: `<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.01//EN" "http://www.w3.org/TR/html4/strict.dtd">`,
	"4strict" 		: `<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.01//EN" "http://www.w3.org/TR/html4/strict.dtd">`,
	"4frameset" 	: `<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.01 Transitional//EN" "http://www.w3.org/TR/html4/loose.dtd">`,
	"4transitional" : `<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.01 Frameset//EN" "http://www.w3.org/TR/html4/frameset.dtd"> `,
}

type DoctypeNode struct {
	NodeType
	Pos
	tr    *Tree
	Doctype string
}

func (t *Tree) newDoctype(pos Pos, dt string) *DoctypeNode {
	return &DoctypeNode{tr: t, NodeType: NodeDoctype, Pos: pos, Doctype: dt}
}

func (d *DoctypeNode) String() string {
	if dt, ok := doctype[d.Doctype]; ok {
		return fmt.Sprintf("%s\n", dt)
	}
	return fmt.Sprintf("<!DOCTYPE html>\n")
}

func (d *DoctypeNode) tree() *Tree {
	return d.tr
}
func (d *DoctypeNode) Copy() Node {
	return &DoctypeNode{tr: d.tr, NodeType: NodeDoctype, Pos: d.Pos, Doctype: d.Doctype}
}
