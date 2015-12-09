package jade

func Parse(name, text string) (*tree, error) {
	return newTree(name).Parse(text, leftDelim, rightDelim, make(map[string]*tree))
}

func (t *tree) String() string {
	return t.Root.String()
}
