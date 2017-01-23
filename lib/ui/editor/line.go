package editor

import (
	"math"
	"unicode/utf8"

	"github.com/gizak/termui"
	rw "github.com/mattn/go-runewidth"
)

type Line struct {
	ContentStartX, ContentStartY int

	Editor *Editor
	Data   []byte
	Cells  []termui.Cell
	Next   *Line
	Prev   *Line
}

func (p *Editor) InitNewLine() *Line {
	p.LinesLocker.Lock()
	defer p.LinesLocker.Unlock()

	ret := &Line{
		Editor:        p,
		ContentStartX: p.Block.InnerArea.Min.X,
		ContentStartY: p.Block.InnerArea.Min.Y,
		Data:          make([]byte, 0),
	}
	p.Lines = append(p.Lines, ret)
	p.DisplayLinesRange[1] = len(p.Lines) - 1

	if nil == p.FirstLine {
		p.FirstLine = ret
	}

	if nil != p.LastLine {
		p.LastLine.Next = ret
		ret.Prev = p.LastLine
	}

	p.LastLine = ret

	p.DisplayLinesRange[1] += 1
	maxHeight := p.InnerArea.Dy()
	maxWidth := p.InnerArea.Dx()
	height := 1
	index := 0
	for i := p.DisplayLinesRange[1] - 2; i >= 0; i-- {
		height += int(math.Ceil(float64(rw.StringWidth(string(p.Lines[i].Data))) / float64(maxWidth)))
		if height > maxHeight {
			break
		}
		index = i
	}
	p.DisplayLinesRange[0] = index

	return ret
}

func (p *Editor) RemoveLine(line *Line) {
	p.LinesLocker.Lock()
	defer p.LinesLocker.Unlock()

	p.CurrentLine = line.Prev

	if nil != line.Prev {
		line.Prev.Next = line.Next
	}
	if nil != line.Next {
		line.Next.Prev = line.Prev
	}

	if p.FirstLine == line {
		p.FirstLine = p.FirstLine.Next
	}

	if p.LastLine == line {
		p.LastLine = p.LastLine.Prev
	}

	for k, v := range p.Lines {
		if line == v {
			p.Lines = append(p.Lines[:k], p.Lines[k+1:]...)
		}
	}
}

func (p *Editor) ClearLines() {
	p.LinesLocker.Lock()
	defer p.LinesLocker.Unlock()

	p.FirstLine = nil
	p.LastLine = nil
	p.CurrentLine = nil
	p.Lines = []*Line{}
	p.CursorLocation.ResetLocation()
	p.DisplayLinesRange = [2]int{0, 1}
}

func (p *Line) Write(ch string) {
	off := p.Editor.CursorLocation.OffXCellIndex

	if off >= len(p.Data) {
		p.Data = append(p.Data, []byte(ch)...)

	} else if 0 == off {
		p.Data = append([]byte(ch), p.Data...)

	} else {
		newData := make([]byte, len(p.Data)+len(ch))
		copy(newData, p.Data[:off])
		copy(newData[off:], []byte(ch))
		copy(newData[off+len(ch):], p.Data[off:])
		p.Data = newData
	}

	fg, bg := p.Editor.TextFgColor, p.Editor.TextBgColor
	cells := termui.DefaultTxBuilder.Build(ch, fg, bg)
	p.Editor.CursorLocation.OffXCellIndex += len(cells)
}

func (p *Line) Backspace() {
	if 0 == len(p.Data) {
		return
	}
	_, rlen := utf8.DecodeLastRune(p.Data)
	p.Data = p.Data[:len(p.Data)-rlen]
}
