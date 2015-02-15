package jade

const (
	tabSize = 4

	outputIndent  = "    "
	nestIndent = true
	lineIndent = false
	
	leftDelim    = "{{"
	rightDelim   = "}}"

	tabComment    = "//-"
	htmlComment   = "//"

	interDelim      = "#{"
	unEscInterDelim = "!{"
	rightInterDelim = "}"

	tabDelim      = "-"
	tabBufDelim   = "="
	tabUnEscDelim = "!="
)



var itemToStr = map[itemType]string {
	itemError:			"itemError",
	itemEOF:			"itemEOF",
	itemEndL:			"itemEndL",
	itemIdentSpace:		"itemIdentSpace",
	itemIdentTab:		"itemIdentTab",
	itemTag:			"itemTag",
	itemVoidTag:		"itemVoidTag",
	itemInlineTag:		"itemInlineTag",
	itemInlineVoidTag:	"itemInlineVoidTag",
	itemHtmlTag:		"itemHtmlTag",
	itemDiv:			"itemDiv",
	itemId:				"itemId",
	itemClass:			"itemClass",
	itemAttr:			"itemAttr",
	itemAction:			"itemAction",
	itemInlineAction:	"itemInlineAction",
	itemInlineText:		"itemInlineText",
	itemFilter:			"itemFilter",
	itemDoctype:		"itemDoctype",
	itemComment:		"itemComment",
	itemBlank:			"itemBlank",
	itemSpace:			"itemSpace",
	itemText:			"itemText",
}


var key = map[string]itemType {
	"area": 	itemVoidTag,
	"base": 	itemVoidTag,
	"col": 		itemVoidTag,
	"command":  itemVoidTag,
	"embed": 	itemVoidTag,
	"hr": 		itemVoidTag,
	"input": 	itemVoidTag,
	"keygen": 	itemVoidTag,
	"link": 	itemVoidTag,
	"meta": 	itemVoidTag,
	"param": 	itemVoidTag,
	"source": 	itemVoidTag,
	"track": 	itemVoidTag,
	"wbr": 		itemVoidTag,

	"include":  itemAction,
	"extends":  itemAction,
	"mixin": 	itemAction,
	"block": 	itemAction,
	"each": 	itemAction,
	"while": 	itemAction,
	"if": 		itemAction,
	"else": 	itemAction,
	"case": 	itemAction,
	"when": 	itemAction,
	"default":  itemAction,

	"a":		itemInlineTag,
	"abbr":		itemInlineTag,
	"acronym":	itemInlineTag,
	"b":		itemInlineTag,
	"code":		itemInlineTag,
	"em":		itemInlineTag,
	"font":		itemInlineTag,
	"i":		itemInlineTag,
	"ins":		itemInlineTag,
	"kbd":		itemInlineTag,
	"map":		itemInlineTag,
	"samp":		itemInlineTag,
	"small":	itemInlineTag,
	"span":		itemInlineTag,
	"strong":	itemInlineTag,
	"sub":		itemInlineTag,
	"sup":		itemInlineTag,

	"br": 		itemInlineVoidTag,
	"img":		itemInlineVoidTag,
}


