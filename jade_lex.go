package jade

import (
	// "fmt"
	"strings"
)



// itemType identifies the type of lex items.
type itemType int
const (
	itemError        itemType = iota // error occurred; value is text of error
	itemEOF
	itemEndL
	itemEndTag
	itemEndAttr

	itemIdentSpace
	itemIdentTab

	itemTag 			// html tag
	itemDiv 			// html div for . or #
	itemInlineTag 		// inline tags
	itemVoidTag 		// self-closing tags
	itemInlineVoidTag 	// inline + self-closing tags
	itemComment

	itemId				// id    attribute
	itemClass			// class attribute
	itemAttr 			// html  attribute value
	itemAttrN			// html  attribute value without quotes
	itemAttrName 		// html  attribute name
	itemAttrVoid 		// html  attribute without value

	itemParentIdent 	// Ident for 'tag:'
	itemText 			// plain text
	itemInlineText
	itemHtmlTag 		// html <tag>

	itemDoctype 		// Doctype tag
	itemBlank
	itemFilter
	itemAction 			// from go template {{...}}
	itemActionEnd 		// from go template {{...}} {{end}}
	itemInlineAction	// title= .titleName

	itemDefine
	itemElse
	itemEnd
	itemIf
	itemRange
	itemNil
	itemTemplate
	itemWith
)



// run runs the state machine for the lexer.
func (l *lexer) run() {
	for l.state = lexDoc; l.state != nil; {
		l.state = l.state(l)
	}
}


func lexDoc(l *lexer) stateFn {
	lexIndents(l)
	var flag = false
	if strings.HasPrefix(l.input[l.pos:], "doctype") {
		l.start += 8
		l.pos = l.start
		l.pos++
		flag = true
	}
	if strings.HasPrefix(l.input[l.pos:], "!!!") {
		l.start += 4
		l.pos = l.start
		l.pos++
		flag = true
	}
	if flag {
		for {
			switch r := l.next(); {
			case isAlphaNumeric(r):
				// absorb.
			default:
				l.backup()
				l.emit(itemDoctype)
				return lexAfterTag
			}
		}
	}
	return lexTags
}

func lexComment(l *lexer) stateFn {
	l.next()
	l.next()
	l.emit(itemComment)
	return lexAfterTag
}

func lexCommentSkip(l *lexer) stateFn {
	l.next()
	l.next()
	l.next()
	l.emit(itemBlank)
	return lexLongText
}

func lexIndents(l *lexer) stateFn {
	l.parenDepth = 0
	for {
		switch l.next() {
		case ' ':
			l.emit(itemIdentSpace)
			l.parenDepth ++
		case '\t':
			l.emit(itemIdentTab)
			l.parenDepth += tabSize
		default:
			l.backup()
			return lexTags
		}
	}
}


// lexTags scans tags.
func lexTags(l *lexer) stateFn {
	if strings.HasPrefix(l.input[l.pos:], l.leftDelim) 	{ return lexAction }
	if strings.HasPrefix(l.input[l.pos:], tabComment) 	{ return lexCommentSkip }
	if strings.HasPrefix(l.input[l.pos:], htmlComment)  { return lexComment }

	switch r := l.next(); {
	case r == eof:
		l.emit(itemEOF)
		// fmt.Println("lex dixi")
		return nil
	case r == '\r':
		return lexTags
	case r == '\n':
		l.emit(itemEndL)
		return lexIndents

	case r == '.':
		l.emit(itemDiv)
		return lexClass
	case r == '#':
		l.emit(itemDiv)
		return lexId
	case r == '|':
		l.ignore()
		return lexTextEndL
	case r == ':':
		l.ignore()
		return lexFilter
	case r == '<':
		return lexHtmlTag
	case r == '+':
		l.ignore()
		if l.toEndL(itemTemplate) { return lexIndents }
		return nil
	case r == '=' || r == '-' || r == '$':
		l.ignore()
		return lexActionEndL

	case isAlphaNumeric(r):
		l.backup()
		return lexTagName

	default:
		return l.errorf("lexTags: %#U", r)
	}
}

func lexAfterTag(l *lexer) stateFn {
	switch r := l.next(); {
	case r == eof:
		l.emit(itemEOF)
		return nil
	case r == '(':
		l.ignore()
		return lexAttr
	case r == '/':
		l.emit(itemEndTag)
		return lexAfterTag
	case r == ':':
		l.ignore()
		lexSp(l)
		return lexTags
	case r == ' ' || r == '\t':
		l.ignore()
		return lexInlineText
	case r == '=':
		l.ignore()
		return lexInlineAction
	case r == '!':
		sp := l.peek()
		l.ignore()
		if sp == '=' { l.next(); l.ignore(); return lexInlineAction }
		return l.errorf("expect '=' after '!'")
	case r == '#':
		l.ignore()
		return lexId
	case r == '.':
		sp := l.peek()
		l.ignore()
		if sp == ' ' { l.next(); l.ignore(); return lexLongText } // { return l.errorf("expect new line after '.'") }
		if sp == '\r' || sp == '\n' { return lexLongText }
		return lexClass
	case r == '\r':
		l.next()
		l.emit(itemEndL)
		return lexIndents
	case r == '\n':
		l.emit(itemEndL)
		return lexIndents
	default:
		return l.errorf("lexAfterTag: %#U", r)
	}
}



func lexId(l *lexer) stateFn {
	var r rune
	for {
		r = l.next()
		if !isAlphaNumeric(r) {
			if r == '#' { return l.errorf("lexId: unexpected #") }
			l.backup()
			break
		}
	}
	l.emit(itemId)
	return lexAfterTag
}

func lexClass(l *lexer) stateFn {
	var r rune
	for {
		r = l.next()
		if !isAlphaNumeric(r) {
			l.backup()
			break
		}
	}
	l.emit(itemClass)
	return lexAfterTag
}

func lexFilter(l *lexer) stateFn {
	var r rune
	for {
		r = l.next()
		if !isAlphaNumeric(r) {
			l.backup()
			break
		}
	}
	l.emit(itemFilter)
	return lexLongText
}

func lexTagName(l *lexer) stateFn {
	for {
		switch r := l.next(); {
		case isAlphaNumeric(r):
			// absorb.
		case r == '-':
			// absorb.	
		default:
			l.backup()
			word := l.input[l.start:l.pos]
			switch key[word] {
			case itemAction:
				if l.toEndL(itemAction)    { return lexIndents }
				return nil
			case itemActionEnd:
				if l.toEndL(itemActionEnd) { return lexIndents }
				return nil
			case itemVoidTag,
				 itemInlineVoidTag,
				 itemInlineTag:
				l.emit(key[word])
			default:
				l.emit(itemTag)
			}
			return lexAfterTag
		}
	}
}


func lexAttr(l *lexer) stateFn {
	for {
		switch r := l.next(); {
		case isAlphaNumeric(r):
			l.backup()
			l.ignore()
			return lexAttrName
		case r == ')':
			l.ignore()
			return lexAfterTag
		case r == ' ' || r == ',' || r == '\t':
		case r == eof:
			return l.errorf("lexAttr: expected ')'")
		}
	}
}
func lexAttrName(l *lexer) stateFn {
	for {
		switch r := l.next(); {
		case isAlphaNumeric(r):
		case r == '-':
			// absorb.
		case r == '=':
			word := l.input[l.start:l.pos]
			switch {
			case word == "id=":
				l.ignore()
				return lexAttrId
			case word == "class=":
				l.ignore()
				return lexAttrClass
			default:
				l.backup(); l.emit(itemAttrName)
				l.next(); l.ignore()
				return lexAttrVal
			}
		case r == ' ' || r == ',' || r == ')' || r == '\r' ||  r == '\n':
			l.backup();
			l.emit(itemAttrVoid)
			return lexAttr
		default:
			return l.errorf("lexAttrName: expected '=' or ' ' %#U", r)
		}
	}
}
func lexAttrId(l *lexer) stateFn {
	stopCh := l.next()
	if stopCh == '"' || stopCh == '\'' {
		l.ignore()
		l.toStopCh(stopCh, itemId, true)
	} else {
		l.toStopSpace(itemId)
	}
	return lexAttr
}
func lexAttrClass(l *lexer) stateFn {
	stopCh := l.next()
	if stopCh == '"' || stopCh == '\'' {
		l.ignore()
		l.toStopCh(stopCh, itemClass, true)
	} else {
		l.toStopSpace(itemClass)
	}
	return lexAttr
}
func lexAttrVal(l *lexer) stateFn {
	stopCh := l.next()
	if stopCh == '"' || stopCh == '\'' {
		l.toStopCh(stopCh, itemAttr, false)
	} else {
		l.toStopSpace(itemAttrN)
	}
	return lexAttr
}
func (l *lexer) toStopCh(stopCh rune, item itemType, backup bool) {
	for {
		switch r := l.next(); {
		case r == stopCh:
			if backup { l.backup() }
			l.emit(item)
			return
		case r == eof || r == '\r' ||  r == '\n':
			l.errorf("toStopCh: expected '%s' %#U",stopCh, r)
			return
		}
	}
}
func (l *lexer) toStopSpace(item itemType) {
	for {
		switch r := l.next(); {
		case r == ' ' || r == ',' || r == ')' || r == '\r' ||  r == '\n':
			l.backup()
			l.emit(item)
			return
		case r == eof:
			l.errorf("toStopCh: expected ')' %#U", r)
			return
		}
	}
}


func lexLongText(l *lexer) stateFn {
	if lexInlineText(l) == nil  { return nil }
	depth := l.toEndIdents()
	for l.parenDepth < depth {
		if lexTextEndL(l) == nil  { return nil }
		depth = l.toEndIdents()
	}
	return lexIndents
}

func (l *lexer) toEndIdents() int {
	var ident int
	for {
		switch l.next() {
		case ' ':
			l.emit(itemIdentSpace)
			ident ++
		case '\t':
			l.emit(itemIdentTab)
			ident += tabSize
		case '\r', '\n':
			ident = 500	// for empty lines in lexLongText
			l.backup()
			l.emit(itemText)
			return ident
		default:
			l.backup()
			return ident
		}
	}
}


func lexHtmlTag(l *lexer) stateFn {
	if l.toEndL(itemHtmlTag) { return lexIndents }
	return nil
}

func lexTextEndL(l *lexer) stateFn {
	if l.toEndL(itemText) { return lexIndents }
	return nil
}

func lexInlineText(l *lexer) stateFn {
	if l.toEndL(itemInlineText) { return lexIndents }
	return nil
}

func lexActionEndL(l *lexer) stateFn {
	if l.toEndL(itemAction) { return lexIndents }
	return nil
}

func lexInlineAction(l *lexer) stateFn {
	if l.toEndL(itemInlineAction) { return lexIndents }
	return nil
}


func lexAction(l *lexer) stateFn {
	l.next()
	l.next()
	l.ignore()
	for {
		l.next()
		if strings.HasPrefix(l.input[l.pos:], l.rightDelim) {
			break
		}
	}
	l.emit(itemAction)
	l.next()
	l.next()
	l.ignore()
	return lexAfterTag
}


func lexSp(l *lexer) {
	for isSpace(l.peek()) {
		l.next()
	}
	l.emit(itemParentIdent)
}


func (l *lexer) toEndL(item itemType) bool {
	Loop:
	for {
		switch r := l.next(); {
		case r == eof:
			if l.pos > l.start { l.emit(item) }
			l.emit(itemEOF)
			return false
		case r == '\r':
			l.backup()
			if l.pos > l.start { l.emit(item) }
			l.next()
			break Loop
		case r == '\n':
			l.backup()
			if l.pos > l.start { l.emit(item) }
			break Loop
		}
	}
	l.next()
	l.emit(itemEndL)
	return true
}
