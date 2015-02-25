package jade

func Parse(name, text string) (*Tree, error) {
	return New(name).Parse(text, leftDelim, rightDelim, make(map[string]*Tree))
}

func (t *Tree) String() string {
	return t.Root.String()
}
