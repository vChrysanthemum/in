package ui

type EditorTmpLinesBuf struct {
	Lines []EditorLine
}

func NewEditorTmpLinesBuf() *EditorTmpLinesBuf {
	ret := &EditorTmpLinesBuf{}
	return ret
}

func (p *EditorTmpLinesBuf) CopyLines(lines []*EditorLine) {
	for _, line := range lines {
		p.Lines = append(p.Lines, *line)
	}
}
