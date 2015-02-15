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
			nest := t.newNest(token.pos, token.val, token.typ, 0, 0)
			t.Root.append( nest )
			if ok, _ := t.parseInside( nest ); !ok { return nil }
		}

		token = t.next()
		// fmt.Printf("%s\t\t\t%s\n", itemToStr[token.typ], token.val)
	}
	fmt.Printf("dixi.")
	return nil
}

func (t *Tree) parseInside( outTag *NestNode ) (bool, int) {
	indentCount := 0
	token := t.next()

	for token.typ != itemEOF {
		// fmt.Printf("-%d-%d---- %s\t\t\t%s\n", indentCount, outTag.Indent, itemToStr[token.typ], token.val)

		switch token.typ {

		case itemEndL:
			indentCount = 0
		case itemIdentSpace:
			indentCount ++
		case itemIdentTab:
			indentCount += tabSize

		case itemText, itemInlineAction:
			outTag.append( t.newLine(token.pos, token.val, token.typ, indentCount, outTag.Nesting + 1) )

		case itemTag, itemDiv, itemInlineTag, itemAction:
			if indentCount > outTag.Indent {
				nest := t.newNest(token.pos, token.val, token.typ, indentCount, outTag.Nesting + 1)

				outTag.append( nest )
				if ok, idt := t.parseInside( nest ); ok {
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
	}
	return false, 0
} 