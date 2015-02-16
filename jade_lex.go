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

	itemIdentSpace
	itemIdentTab

	itemTag 			// html tag
	itemDiv 			// html div for . or #
	itemInlineTag 		// inline tags
	itemVoidTag 		// self-closing tags
	itemInlineVoidTag 	// inline + self-closing tags

	itemId				// id    attribute
	itemClass			// class attribute
	itemAttr 			// html  attribute

	itemParentIdent 	// Ident for 'tag:'
	itemText 			// plain text
	itemInlineText
	itemHtmlTag 		// html <tag>

	itemDoctype 		// Doctype tag
	itemComment
	itemBlank
	itemFilter
	itemAction 			// from go template {{...}}
	itemInlineAction	// title= .titleName
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
	for {
		switch l.next() {
		case ' ':
			l.emit(itemIdentSpace)
		case '\t':
			l.emit(itemIdentTab)
		default:
			l.backup()
			return lexTags
		}
	}
}


// lexTags scans tags.
func lexTags(l *lexer) stateFn {
		// fmt.Println("------")
	if strings.HasPrefix(l.input[l.pos:], l.leftDelim) 	{ return lexAction }
	if strings.HasPrefix(l.input[l.pos:], tabComment) 	{ return lexCommentSkip }
	if strings.HasPrefix(l.input[l.pos:], htmlComment)  { return lexComment }

	switch r := l.next(); {
	case r == eof:
		l.emit(itemEOF)
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
	case r == '=' || r == '+' || r == '-':
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
	case r == '(':
		l.ignore()
		return lexAttr
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
	case r == '#':
		l.ignore()
		return lexId
	case r == '.':
		sp := l.peek()
		l.ignore()
		if sp == ' ' { l.next(); l.ignore(); return lexLongText }
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
		default:
			l.backup()
			word := l.input[l.start:l.pos]
			switch {
			case key[word] >= itemAction:
				return lexActionEndL
			case key[word] > itemTag :
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
		case r == ')':
			l.backup()
			l.emit(itemAttr)
			l.next()
			l.ignore()
			return lexAfterTag
		case r == ',':
			l.emit(itemAttr)
			l.ignore()
		case r == '\n':
			return l.errorf("lexId: expected ')'")
		}
	}
}


func lexLongText(l *lexer) stateFn {
	var (
	 	startIdent int
		nextIdent  int = 1000
	)
	if lexTextEndL(l) == nil  { return nil }
	startIdent = lexSpace(l)
	for startIdent <= nextIdent {
		if lexTextEndL(l) == nil  { return nil }
		nextIdent = lexSpace(l)
	}
	return lexIndents
}


func lexSpace(l *lexer) int {
	var ident int
Loop:
	for {
		switch l.next() {
		case ' ':
			l.emit(itemIdentSpace)
			ident ++
		case '\t':
			l.emit(itemIdentTab)
			ident += tabSize
		default:
			l.backup()
			break Loop
		}
	}
	return ident
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