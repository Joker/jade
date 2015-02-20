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
		case itemHtmlTag:
			t.Root.append( t.newLine(token.pos, token.val, token.typ, 0, 0) )

		case itemTag, itemDiv, itemInlineTag, itemAction, itemComment, itemBlank:
			nest := t.newNest(token.pos, token.val, token.typ, 0, 0)
			if token.typ < itemComment { t.parseAttr( nest ) }
			t.Root.append( nest )
			t.parseInside( nest )
		}

		token = t.next()
		// fmt.Printf("%s\t\t\t%s\n", itemToStr[token.typ], token.val)
	}
	fmt.Println("dixi.")
	return nil
}

func (t *Tree) parseInside( outTag *NestNode ) int {
	indentCount := 0
	token := t.next()

	for token.typ != itemEOF {
		// fmt.Printf("-%d-%d---- %s\t\t\t%s\n", indentCount, outTag.Indent, itemToStr[token.typ], token.val)

		switch token.typ {
		case itemError:
			t.errorf("%s", token.val)

		case itemEndL:
			indentCount = 0
		case itemIdentSpace:
			indentCount ++
		case itemIdentTab:
			indentCount += tabSize
		case itemParentIdent:
			indentCount = outTag.Indent + 1  // for  "tag: tag: tag"

		case itemText, itemInlineText, itemInlineAction:
			outTag.append( t.newLine(token.pos, token.val, token.typ, indentCount, outTag.Nesting + 1) )

		case itemHtmlTag:
			if indentCount > outTag.Indent {
				outTag.append( t.newLine(token.pos, token.val, token.typ, indentCount, outTag.Nesting + 1) )
			}else{
				t.backup(); return indentCount
			}

		case itemVoidTag, itemInlineVoidTag:
			if indentCount > outTag.Indent {
				nest := t.newNest(token.pos, token.val, token.typ, indentCount, outTag.Nesting + 1)
				t.parseAttr( nest )
				outTag.append( nest )
			}else{
				t.backup(); return indentCount
			}

		case itemTag, itemDiv, itemInlineTag:
			if indentCount > outTag.Indent {
				nest := t.newNest(token.pos, token.val, token.typ, indentCount, outTag.Nesting + 1)
				t.parseAttr( nest )
				outTag.append( nest )
				indentCount = t.parseInside( nest )
			}else{
				t.backup(); return indentCount
			}

		case itemAction:
			if indentCount > outTag.Indent {
				nest := t.newNest(token.pos, token.val, token.typ, indentCount, outTag.Nesting + 1)
				outTag.append( nest )
				indentCount = t.parseInside( nest )
			}else{
				t.backup(); return indentCount
			}

		case itemBlank, itemComment:
			if indentCount > outTag.Indent {
				nest := t.newNest(token.pos, token.val, token.typ, indentCount, outTag.Nesting + 1)
				if token.typ == itemComment { outTag.append( nest ) }
				indentCount = t.parseInside( nest )
			}else{
				t.backup(); return indentCount
			}

		}
		token = t.next()
	}
	t.backup(); return indentCount
}

func (t *Tree) parseAttr( currentTag *NestNode ) {
	for {
		attr := t.next()
		// fmt.Println(itemToStr[attr.typ], attr.val)
		switch attr.typ {
		case itemError:
			t.errorf("%s", attr.val)
		case itemId:
			if len(currentTag.id) > 0 { t.errorf("unexpected second id \"%s\" ", attr.val) }
			currentTag.id = attr.val
		case itemClass:
			currentTag.class = append( currentTag.class, attr.val )
		case itemAttr, itemAttrN, itemAttrName, itemAttrVoid:
			currentTag.append( t.newAttr(attr.pos, attr.val, attr.typ) )
		default:
			t.backup()
			currentTag.append( t.newAttr(attr.pos, ">", itemEndAttr) )
			return
		}
	}
}