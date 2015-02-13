package jade

import (
	"fmt"
)



func (t *Tree) parse(treeSet map[string]*Tree) (next Node) {
	token := t.next()
	t.Root = t.newList(token.pos)

	for token.typ != itemEOF {

		switch token.typ {
		case itemError:
			t.errorf("%s", token.val)
		case itemDoctype:
			t.Root.append( t.newDoctype(token.pos, token.val) )
		case itemTag:
			tag := t.newTag(token.pos, token.val)
			t.Root.append( tag )
			if ok, _ := t.parseInside( tag ); !ok { return nil }
		}

		token = t.next()
		// fmt.Printf("%s\t\t\t%s\n", itemToStr[token.typ], token.val)
	}
	fmt.Printf("dixi.")
	return nil
}

func (t *Tree) parseInside( outTag *TagNode ) (bool, int) {
	indentCount := 0
	token := t.next()

	for token.typ != itemEOF {

		switch token.typ {
		case itemEndL:
			indentCount = 0
		case itemIdentSpace:
			indentCount ++
		case itemIdentTab:
			indentCount += tabSize
		case itemText:
			outTag.append( t.newText(token.pos, token.val) )
		case itemTag:
			if indentCount > outTag.Indent {
				tag := t.newTag(token.pos, token.val)
				tag.Indent  = indentCount
				tag.Nesting = outTag.Nesting + 1
				outTag.append( tag )
				if ok, idt := t.parseInside( tag ); ok {
					indentCount = idt
				} else {
					return false, 0
				}
			}else{
				t.backup()
				return true, indentCount
			}
		}

		token = t.next()
		// fmt.Printf("-%d-%d---- %s\t\t\t%s\n", indentCount, outTag.Indent, itemToStr[token.typ], token.val)
	}
	return false, 0
}
